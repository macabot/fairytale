package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/macabot/fairytale/internal/model"
	"github.com/macabot/fairytale/internal/set"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type moduleList struct {
	Dir   string
	GoMod string
}

const fairytaleModPath = "github.com/macabot/fairytale"

var hub *Hub

var (
	serverPrefix  = "server"
	watcherPrefix = "watcher"
	maxPrefixLen  = len(watcherPrefix)
)

func serverLogf(format string, a ...any) {
	logf(ServerLog, format, a...)
}

func watcherLogf(format string, a ...any) {
	logf(WatcherLog, format, a...)
}

type LogKind int

const (
	ServerLog LogKind = iota + 1
	WatcherLog
)

func logf(kind LogKind, format string, a ...any) {
	var c *color.Color
	var prefix string
	switch kind {
	case ServerLog:
		c = color.New(color.FgCyan)
		prefix = serverPrefix
	case WatcherLog:
		c = color.New(color.FgMagenta)
		prefix = watcherPrefix
	default:
		panic("invalid LogKind")
	}
	c.Printf("%-*s | ", maxPrefixLen, prefix)
	fmt.Printf(format, a...)
}

func handle(mainWasmPath, wasmExecJsPath, fairytaleAssetsDir, userAssetsDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var servePath string
		switch r.URL.Path {
		case "/main.wasm":
			servePath = mainWasmPath
		case "/wasm_exec.js":
			servePath = wasmExecJsPath
		default:
			assetFound := false
			if userAssetsDir != "" && r.URL.Path != "/" {
				servePath = filepath.Join(userAssetsDir, filepath.FromSlash(r.URL.Path))
				if _, err := os.Stat(servePath); err == nil {
					assetFound = true
				}
			}
			if !assetFound {
				servePath = filepath.Join(fairytaleAssetsDir, filepath.FromSlash(r.URL.Path))
			}
		}
		servePath = filepath.Clean(servePath)
		http.ServeFile(w, r, servePath)
		serverLogf("[%s] %s %s --> %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL, servePath)
	}
}

func createReloadBytes() []byte {
	message := model.SocketMessage{
		Type: model.SocketMessageReload,
	}
	b, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	return b
}

var reloadBytes = createReloadBytes()

func handleReload(w http.ResponseWriter, r *http.Request) {
	// TODO only allow POST
	hub.broadcast <- reloadBytes
}

func findFairytaleAssetsDir(cmd *cobra.Command) string {
	out, err := exec.Command("go", "list", "-json", "-m").Output()
	cobra.CheckErr(err)
	var mod moduleList
	cobra.CheckErr(json.Unmarshal(out, &mod))

	goMod, err := os.ReadFile(mod.GoMod)
	cobra.CheckErr(err)
	f, err := modfile.Parse(mod.GoMod, goMod, nil)
	cobra.CheckErr(err)

	var modVer *module.Version
	for _, require := range f.Require {
		if require.Mod.Path == fairytaleModPath {
			if require.Indirect {
				cmd.PrintErrf("Warning: '%s' should not be an indirect requirement in your go.mod file.", fairytaleModPath)
			}
			modVer = &require.Mod
			break
		}
	}
	for _, replace := range f.Replace {
		if replace.Old.Path == fairytaleModPath {
			modVer = &replace.New
			break
		}
	}
	if modVer == nil {
		cobra.CheckErr("Your go.mod file must require " + fairytaleModPath)
	}

	var fairytaleDir string
	if modfile.IsDirectoryPath(modVer.Path) {
		if filepath.IsAbs(modVer.Path) {
			fairytaleDir = modVer.Path
		} else {
			fairytaleDir = filepath.Join(mod.Dir, modVer.Path)
		}
	} else {
		modCacheBytes, err := exec.Command("go", "env", "GOMODCACHE").Output()
		cobra.CheckErr(err)
		modCache := strings.TrimSuffix(string(modCacheBytes), "\n")
		path := append([]string{modCache}, strings.Split(modVer.String(), "/")...)
		fmt.Println("modCache", path)
		fairytaleDir = filepath.Join(path...)
	}
	return filepath.Join(fairytaleDir, "cmd", "fairytale", "cmd", "assets")
}

func findWasmExecJsPath() string {
	if path := os.Getenv("WASM_EXEC_JS_PATH"); path != "" {
		return path
	}
	out, err := exec.Command("go", "env", "GOROOT").Output()
	cobra.CheckErr(err)
	goRoot := strings.TrimSuffix(string(out), "\n")
	return filepath.Join(goRoot, "lib", "wasm", "wasm_exec.js")
}

func runWatcher(stop chan struct{}, hub *Hub, paths []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for i := range paths {
		paths[i] = filepath.Clean(paths[i])
	}

	dirPaths := set.Set[string]{}
	upgradedDirPaths := map[string]set.Set[string]{}
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			dirPaths.Add(path)
		} else {
			dir := filepath.Dir(path)
			if upgradedPaths, ok := upgradedDirPaths[dir]; ok {
				upgradedPaths.Add(path)
			} else {
				upgradedDirPaths[dir] = set.New(path)
			}
		}
	}
	for dir := range dirPaths {
		delete(upgradedDirPaths, dir)
	}

	go func() {
		for {
			select {
			// FIXME why doesn't VSC recognize watcher.Events and watcher.Errors?
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				name := filepath.Clean(event.Name)

				// Ignore events of files that are only watched because we are watching the parent of another file we're interested in.
				dir := filepath.Dir(name)
				if paths, ok := upgradedDirPaths[dir]; ok {
					if !paths.Has(name) {
						continue
					}
				}

				watcherLogf("[%s] event %v\n", time.Now().Format(time.RFC3339), event)
				hub.broadcast <- reloadBytes
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				watcherLogf("[%s] error: %v\n", time.Now().Format(time.RFC3339), err)
			}
		}
	}()

	pathsToWatch := dirPaths.Clone()
	for dir := range upgradedDirPaths {
		pathsToWatch.Add(dir)
	}
	for path := range dirPaths {
		watcherLogf("[%s] watch %s\n", time.Now().Format(time.RFC3339), path)
		if err := watcher.Add(path); err != nil {
			log.Fatal(err)
		}
	}
	for dir, paths := range upgradedDirPaths {
		for path := range paths {
			watcherLogf("[%s] watch %s\n", time.Now().Format(time.RFC3339), path)
		}
		// Always watch directories in case a file is deleted and recreated.
		if err := watcher.Add(dir); err != nil {
			log.Fatal(err)
		}
	}

	<-stop
}

var (
	watch         bool
	userAssetsDir string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve address main.wasm",
	Short: "Serve the fairytale application",
	Long: `Serve the fairytale application during development of your applications. The fairytale app will be served on the given address.

Use the WASM_EXEC_JS_PATH environment variable to override the path to the wasm_exec.js file.`,
	Example: "fairytale serve :8080 path/to/main.wasm",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		mainWasmPath := args[1]

		if userAssetsDir != "" {
			if _, err := os.Stat(userAssetsDir); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					cobra.CheckErr(fmt.Errorf("assets directory '%s' does not exist", userAssetsDir))
				} else {
					cobra.CheckErr(fmt.Errorf("could not check if assets directory '%s' exist: %w", userAssetsDir, err))
				}
			}
		}
		fairytaleAssetsDir := findFairytaleAssetsDir(cmd)
		wasmExecJsPath := findWasmExecJsPath()

		hub = newHub()
		go hub.run()

		if watch {
			stopWatcher := make(chan struct{})
			defer func() {
				stopWatcher <- struct{}{}
			}()
			paths := []string{mainWasmPath, wasmExecJsPath, fairytaleAssetsDir}
			if userAssetsDir != "" {
				paths = append(paths, userAssetsDir)
			}
			go runWatcher(stopWatcher, hub, paths)
		}

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			serveWs(hub, w, r)
		})
		http.HandleFunc("/reload", handleReload)
		http.HandleFunc("/", handle(mainWasmPath, wasmExecJsPath, fairytaleAssetsDir, userAssetsDir))
		cobra.CheckErr(http.ListenAndServe(address, nil))
	},
}

func init() {
	serveCmd.Flags().BoolVar(&watch, "watch", false, "Watch for changes made to the files and reload the page when it happens.")
	serveCmd.Flags().StringVar(&userAssetsDir, "assets", "", "Serve the files in the given directory.")
	rootCmd.AddCommand(serveCmd)
}
