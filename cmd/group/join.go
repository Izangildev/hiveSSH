package group

import (
	"fmt"
	"hivessh/internal/store"

	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join <group> <server_id>",
	Short: "Add a server to a group",
	Long: `Add a stored server to an existing group.

Usage examples:
  hivessh group join <group_name> <server_id>`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		groupname := args[0]
		serverID := args[1]

		if err := store.JoinServer(groupname, serverID); err != nil {
			fmt.Printf("[❌] Failed to add server to group: %s\n", err)
			return
		}
		fmt.Printf("[✅] Server '%s' successfully added to group '%s'\n", serverID, groupname)
	},
}

func init() {
	groupCmd.AddCommand(joinCmd)
}
