package cmd

import (
	"hivessh/logic"
	"os"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Command run, used to execute a unique remote command.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var command = os.Args[2]
		logic.Run(command)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
