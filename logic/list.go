package logic

import (
	"fmt"
)

func List() {
	// Header
	fmt.Printf("%-10s %-18s %-14s\n", "NAME", "IP", "SSH STATUS")
	fmt.Println("────────── ────────────────── ──────────────")

	for name, ip := range servers {
		ping := getStatus(ip)
		var status string

		switch ping {
		case true:
			status = "reachable"
		case false:
			status = "unreachable"
		}

		fmt.Printf("%-10s %-18s %-14s\n", name, ip, status)
		fmt.Println("────────── ────────────────── ──────────────")
	}
}
