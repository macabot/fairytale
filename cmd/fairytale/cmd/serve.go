package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/macabot/fairytale/internal/model"
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

func handle(mainWasmPath, wasmExecJsPath, assetsDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var servePath string
		switch r.URL.Path {
		case "/main.wasm":
			servePath = mainWasmPath
		case "/wasm_exec.js":
			servePath = wasmExecJsPath
		default:
			servePath = filepath.Join(assetsDir, filepath.FromSlash(r.URL.Path))
			if r.URL.Path == "/" {
				servePath = filepath.Join(servePath, "index.html")
			}
		}
		servePath = filepath.Clean(servePath)
		http.ServeFile(w, r, servePath)
		fmt.Printf("[%s] %s %s --> %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL, servePath)
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

func findAssetsDir(cmd *cobra.Command) string {
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
				cmd.PrintErrf("Warning: '%s' should not be an indirect requirement in your go.mod file.")
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
		cobra.CheckErr("Your go.mod file does not require " + fairytaleModPath)
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
	out, err := exec.Command("go", "env", "GOROOT").Output()
	cobra.CheckErr(err)
	goRoot := strings.TrimSuffix(string(out), "\n")
	return filepath.Join(goRoot, "misc", "wasm", "wasm_exec.js")
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve [address] [main.wasm]",
	Short:   "Serve the fairytale application",
	Long:    `Serve the fairytale application during development of your applications. The fairytale app will be served on the given address.`,
	Example: "fairytale serve :8080 path/to/main.wasm",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		mainWasmPath := args[1]

		assetsDir := findAssetsDir(cmd)
		wasmExecJsPath := findWasmExecJsPath()

		hub = newHub()
		go hub.run()

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			serveWs(hub, w, r)
		})
		http.HandleFunc("/reload", handleReload)
		http.HandleFunc("/", handle(mainWasmPath, wasmExecJsPath, assetsDir))
		cobra.CheckErr(http.ListenAndServe(address, nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
