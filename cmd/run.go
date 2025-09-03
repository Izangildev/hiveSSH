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

		if command == "" {
			fmt.Println("[✖] Command cannot be empty. Please provide a valid command to run.")
			return
		}

		if target == "" {
			fmt.Println("[✖] You must specify a target server using --to")
			return
		}

		fmt.Printf("[→] Connecting to server '%s'...\n", target)

		if err := logic.Run(command, target); err != nil {
			fmt.Printf("[✖] Command execution failed: %s\n", err)
			return
		}

		fmt.Printf("[✔] Command executed successfully on '%s'\n", target)
	},
}

func init() {
	runCmd.Flags().StringVar(&target, "to", "", "IP or name of the server stored in DB")
	rootCmd.AddCommand(runCmd)
}
