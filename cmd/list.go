package cmd

import (
	"hivessh/logic"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display the configured hosts and their status",
	Long: `The list command shows all hosts registered in HiveSSH 
along with their IP address and connection status (reachable or unreachable).`,
	Run: func(cmd *cobra.Command, args []string) {

		logic.List()
	},
}

func init() {
	runCmd.Flags().StringVar(&target, "to", "", "IP or name of the server stored in DB")
	rootCmd.AddCommand(listCmd)
}
