package main

import (
	"embed"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
)

//go:embed public/*
var public embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Not enough arguments.")
		os.Exit(1)
	}
	if os.Args[1] != "init" {
		fmt.Fprintf(os.Stderr, "Invalid argument '%s'.\n", os.Args[1])
		os.Exit(1)
	}
	writeDir := "."
	if len(os.Args) == 3 {
		writeDir = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "Too many arguments.")
		os.Exit(1)
	}
	if _, err := os.Stat(writeDir); os.IsNotExist(err) {
		if err := os.MkdirAll(writeDir, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Could not create directory '%s': %s\n", writeDir, err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Could not check if directory '%s' exists: %s\n", writeDir, err)
		os.Exit(1)
	}

	publicDir, err := public.ReadDir("public")
	if err != nil {
		panic(fmt.Errorf("fairytale: could not read embedded public directory: %w", err))
	}
	for _, fileEntry := range publicDir {
		readPath := filepath.Join("public", fileEntry.Name())
		b, err := public.ReadFile(readPath)
		if err != nil {
			panic(fmt.Errorf("fairytale: could not read embedded file '%s': %w", readPath, err))
		}

		writePath := filepath.Join(writeDir, fileEntry.Name())
		if err := os.WriteFile(writePath, b, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Could not write file '%s': %s\n", writePath, err)
			os.Exit(1)
		}
	}

	goRoot := os.Getenv("GOROOT")
	if goRoot == "" {
		goRoot = build.Default.GOROOT
	}
	wasmExecReadPath := filepath.Join(goRoot, "misc", "wasm", "wasm_exec.js")
	wasmExec, err := ioutil.ReadFile(wasmExecReadPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file '%s': %s\n", wasmExecReadPath, err)
		os.Exit(1)
	}

	wasmExecWritePath := filepath.Join(writeDir, "wasm_exec.js")
	if err := os.WriteFile(wasmExecWritePath, wasmExec, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Could not write file '%s': %s\n", wasmExecWritePath, err)
		os.Exit(1)
	}
}
