package cmd

import (
	"fmt"
	"hivessh/logic"

	"github.com/spf13/cobra"
)

var outputType string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display the configured hosts and their status",
	Long: `The list command shows all hosts registered in HiveSSH 
along with their IP address and connection status (reachable or unreachable).`,
	Run: func(cmd *cobra.Command, args []string) {

		if outputType != "" && outputType != "json" && outputType != "csv" {
			fmt.Println("[❌] Invalid output type. Supported types are: json, csv")
			return
		}

		if err := logic.List(outputType); err != nil {
			fmt.Printf("[❌] Command execution failed: %s\n", err)
			return
		}
	},
}

func init() {
	listCmd.Flags().StringVar(&outputType, "output", "", "Specify output format: json or csv")
	RootCmd.AddCommand(listCmd)
}
