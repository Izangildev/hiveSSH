/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"hivessh/logic"
	"net"

	"github.com/spf13/cobra"
)

var user string
var port int
var description string

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join command. Used to join servers into the database.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println("[❌] You must provide both a server name and an IP address.")
			return
		}

		var name = args[0]

		if name == "" {
			fmt.Println("[❌] Server name cannot be empty. Please provide a valid name.")
			return
		}

		var ip = args[1]

		if net.ParseIP(ip) == nil || ip == "" {
			fmt.Println("[❌] Invalid IP address. Please provide a valid IPv4 or IPv6 address.")
			return
		}

		if err := logic.Join(name, ip, user, description, port); err != nil {
			fmt.Printf("[❌] Failed to add server: %s\n", err)
			return
		}

		fmt.Printf("[✅] Server '%s' successfully added with IP '%s'\n", name, ip)

	},
}

func init() {
	joinCmd.Flags().StringVarP(&user, "user", "u", "root", "SSH user for the server (default: root)")
	joinCmd.Flags().IntVarP(&port, "port", "p", 22, "SSH port for the server (default: 22)")
	joinCmd.Flags().StringVarP(&description, "description", "d", "", "Description for the server")
	RootCmd.AddCommand(joinCmd)
}
