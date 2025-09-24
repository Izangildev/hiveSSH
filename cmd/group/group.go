package group

import (
	"fmt"

	"hivessh/cmd"

	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage server groups",
	Long:  `Manage server groups for organizing hosts and executing commands on multiple servers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("[❌] You must specify an action.")
	},
}

func init() {
	cmd.RootCmd.AddCommand(groupCmd)
}
