package cmd

import (
	"fmt"
	"hivessh/internal/ssh"

	"github.com/spf13/cobra"
)

var target string
var targetGroup string

var runCmd = &cobra.Command{
	Use:   "run <command>",
	Short: "Execute a remote command on a server or group",
	Long:  `Run use: hivessh run <command> --to <target> OR hivessh run <command> --group <groupname>`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]

		if target == "" && targetGroup == "" {
			fmt.Println("[❌] You must specify a target using --to <server> or --group <groupname>")
			return
		}
		if target != "" && targetGroup != "" {
			fmt.Println("[❌] You cannot use --to and --group at the same time.")
			return
		}

		if target != "" {
			fmt.Printf("[➡️] Connecting to server '%s'...\n", target)
			if err := ssh.Run(command, target); err != nil {
				fmt.Printf("[❌] Command execution failed: %s\n", err)
			}
		} else {
			fmt.Printf("[➡️] Connecting to group '%s'...\n", targetGroup)
			fmt.Print("------------------------------------------\n")
			if err := ssh.RunGroup(command, targetGroup); err != nil {
				fmt.Printf("[❌] Command execution failed: %s\n", err)
			}
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&target, "to", "", "Name, IP or ID of the target server")
	runCmd.Flags().StringVarP(&targetGroup, "group", "g", "", "Name of the target group")
	RootCmd.AddCommand(runCmd)
}
