package cmd

import (
	"fmt"
	"hivessh/internal/store"
	"strings"

	"github.com/spf13/cobra"
)

var outputType string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display the configured hosts and their status",
	Long:  `Shows all hosts registered in HiveSSH along with their IP address and SSH connection status.`,
	Run: func(cmd *cobra.Command, args []string) {
		out := strings.ToLower(outputType)
		if out != "" && out != "json" && out != "csv" {
			fmt.Println("[❌] Invalid output type. Supported types are: json, csv")
			return
		}
		if err := store.List(out); err != nil {
			fmt.Printf("[❌] Command execution failed: %s\n", err)
		}
	},
}

func init() {
	listCmd.Flags().StringVar(&outputType, "output", "", "Specify output format: json or csv")
	RootCmd.AddCommand(listCmd)
}
