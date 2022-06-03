package cmd

import (
	"embed"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed public/*
var public embed.FS

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initiate a fairytale directory",
	Long: `Initiate a directory with the assets needed to serve the fairytale.
It will initiate the current directory if no directory is provided. It will create a directory if the provided directory doesn't exist.

Examples:
  # Initiate the current directory
  $ fairytale init

  # Initiate the directory named my-fairytale-app
  $ fairytale init my-fairytale-app

You must add the WASM build of your fairytale app, named main.wasm, to the initiated directory. E.g.:
  GOOS=js GOARCH=wasm go build -o path/to/initiated-directory/main.wasm main.go
`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		writeDir := "."
		if len(args) >= 1 {
			writeDir = args[0]
		}
		if _, err := os.Stat(writeDir); os.IsNotExist(err) {
			if err := os.MkdirAll(writeDir, os.ModePerm); err != nil {
				return fmt.Errorf("fairytale: could not create directory '%s': %w", writeDir, err)
			}
		} else if err != nil {
			return fmt.Errorf("fairytale: could not check if directory '%s' exists: %w", writeDir, err)
		}

		publicDir, err := public.ReadDir("public")
		if err != nil {
			return fmt.Errorf("fairytale: could not read embedded public directory: %w", err)
		}
		for _, fileEntry := range publicDir {
			readPath := filepath.Join("public", fileEntry.Name())
			b, err := public.ReadFile(readPath)
			if err != nil {
				return fmt.Errorf("fairytale: could not read embedded file '%s': %w", readPath, err)
			}

			writePath := filepath.Join(writeDir, fileEntry.Name())
			if err := os.WriteFile(writePath, b, 0644); err != nil {
				return fmt.Errorf("fairytale: could not write file '%s': %w", writePath, err)
			}
		}

		goRoot := os.Getenv("GOROOT")
		if goRoot == "" {
			goRoot = build.Default.GOROOT
		}
		wasmExecReadPath := filepath.Join(goRoot, "misc", "wasm", "wasm_exec.js")
		wasmExec, err := ioutil.ReadFile(wasmExecReadPath)
		if err != nil {
			return fmt.Errorf("fairytale: could not read file '%s': %w", wasmExecReadPath, err)
		}

		wasmExecWritePath := filepath.Join(writeDir, "wasm_exec.js")
		if err := os.WriteFile(wasmExecWritePath, wasmExec, 0644); err != nil {
			return fmt.Errorf("fairytale: could not write file '%s': %w", wasmExecWritePath, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
