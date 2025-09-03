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

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join command. Used to join servers into the database.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var name = args[0]

		if name == "" {
			fmt.Println("Error: name cannot be empty")
			return
		}

		var ip = args[1]

		if net.ParseIP(ip) == nil || ip == "" {
			fmt.Println("Error: invalid IP address format")
			return
		}

		if err := logic.Join(name, ip); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
