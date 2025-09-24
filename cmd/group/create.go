package group

import (
	"fmt"
	"hivessh/logic/group"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new server group",
	Long:  `The create command allows you to create a new server group`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("[❌] You must provide a group name.")
			return
		}

		groupName := args[0]
		if groupName == "" {
			fmt.Println("[❌] Group name cannot be empty.")
			return
		}

		if err := group.Create(groupName); err != nil {
			fmt.Printf("[❌] Failed to create group: %s\n", err)
			return
		}

		fmt.Printf("[✅] Group '%s' created successfully.\n", groupName)
	},
}

func init() {
	groupCmd.AddCommand(createCmd)
}
