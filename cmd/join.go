package cmd

import (
	"fmt"
	"hivessh/internal/store"
	"net"

	"github.com/spf13/cobra"
)

var user string
var port int
var description string

var joinCmd = &cobra.Command{
	Use:   "join <name> <ip>",
	Short: "Add a server to the database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		ip := args[1]

		if net.ParseIP(ip) == nil {
			fmt.Println("[❌] Invalid IP address. Please provide a valid IPv4 or IPv6 address.")
			return
		}
		if port < 1 || port > 65535 {
			fmt.Println("[❌] Invalid port. Must be between 1 and 65535.")
			return
		}
		if user == "" {
			fmt.Println("[❌] User cannot be empty.")
			return
		}

		if err := store.Join(name, ip, user, description, port); err != nil {
			fmt.Printf("[❌] Failed to add server: %s\n", err)
			return
		}
		fmt.Printf("[✅] Server '%s' successfully added with IP '%s'\n", name, ip)
	},
}

func init() {
	joinCmd.Flags().StringVarP(&user, "user", "u", "root", "SSH user for the server")
	joinCmd.Flags().IntVarP(&port, "port", "p", 22, "SSH port for the server")
	joinCmd.Flags().StringVarP(&description, "description", "d", "", "Description for the server")
	RootCmd.AddCommand(joinCmd)
}
