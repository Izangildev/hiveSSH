package cmd

import (
	"fmt"
	"hivessh/logic"

	"github.com/spf13/cobra"
)

var target string
var targetGroup string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Command run, used to execute a unique remote command.",
	Long:  `Run use: hivessh run <command> --to <target> OR hivessh run <command> --group <groupname>`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("[❌] You must provide a command to run.")
			return
		}

		var command = args[0]

		if command == "" {
			fmt.Println("[❌] Command cannot be empty. Please provide a valid command to run.")
			return
		}

		if target == "" && targetGroup == "" {
			fmt.Println("[❌] You must specify a target server using --to or a target group using --group")
			return
		}

		if target != "" {
			fmt.Printf("[➡️] Connecting to server '%s'...\n", target)
		} else {
			fmt.Printf("[➡️] Connecting to group '%s'...\n", targetGroup)
		}

		// Execute the command on the specified target
		if target != "" {
			if err := logic.Run(command, target); err != nil {
				fmt.Printf("[❌] Command execution failed: %s\n", err)
				return
			}
		} else {
			if err := logic.RunGroup(command, targetGroup); err != nil {
				fmt.Printf("[❌] Command execution failed: %s\n", err)
				return
			}
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&target, "to", "", "IP or name of the server stored in DB")
	runCmd.Flags().StringVar(&targetGroup, "group", "", "Name of the group to run the command on")
	RootCmd.AddCommand(runCmd)
}
