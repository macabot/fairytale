package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve [address] [directory]",
	Short: "Serve the fairytale application",
	Long: `Serve the fairytale application on the given address. It will use the current directory if no directory is provided.

Examples:
  # Serve the current directory on :8080
  $ fairytale serve :8080

  # Serve the directory named my-fairytale-app on :8080
  $ fairytale serve :8080 my-fairytale-app`,
	Args: cobra.MatchAll(
		cobra.MinimumNArgs(1),
		cobra.MaximumNArgs(2),
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]
		dir := "."
		if len(args) >= 2 {
			dir = args[1]
		}
		http.Handle("/", http.FileServer(http.Dir(dir)))
		return http.ListenAndServe(address, nil)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
