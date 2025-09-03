package cmd

import (
	"fmt"
	"hivessh/logic"

	"github.com/spf13/cobra"
)

var target string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Command run, used to execute a unique remote command.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var command = args[0]

		if target == "" {
			fmt.Println("You need to specify a target (Name or IP).")
			return
		}

		logic.Run(command, target)
	},
}

func init() {
	runCmd.Flags().StringVar(&target, "to", "", "IP or name of the server stored in DB")
	rootCmd.AddCommand(runCmd)
}
