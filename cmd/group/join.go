package group

import (
	"fmt"

	"hivessh/logic/group"

	"github.com/spf13/cobra"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Add servers to a group",
	Long: `The join command allows you to add your stored servers to an existing group.

Usage examples:
  hivessh group join <group_name> <server_id>

This command will associate the specified server with the given group.`,
	Run: func(cmd *cobra.Command, args []string) {
		groupname := args[0]
		serverID := args[1]

		if err := group.JoinServer(groupname, serverID); err != nil {
			fmt.Printf("[❌] Failed to add server to group: %s\n", err)
			return
		}

		fmt.Printf("[✅] Server '%s' successfully added to group '%s'\n", serverID, groupname)
	},
}

func init() {
	groupCmd.AddCommand(joinCmd)
}
