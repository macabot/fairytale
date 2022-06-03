package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fairytale",
	Short: "A CLI for the fairytale library",
	Long:  `Use the CLI to initiate and serve a fairytale application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
