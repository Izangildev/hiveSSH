package group

import (
	"fmt"
	"hivessh/internal/store"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new server group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		if err := store.Create(groupName); err != nil {
			fmt.Printf("[❌] Failed to create group: %s\n", err)
			return
		}
		fmt.Printf("[✅] Group '%s' created successfully.\n", groupName)
	},
}

func init() {
	groupCmd.AddCommand(createCmd)
}
