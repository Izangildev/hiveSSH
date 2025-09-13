package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Class to posteriorly extract server information in json or csv format
type extractableServer struct {
	Name   string
	IP     string
	Status string
}

func List(outputType string) error {
	// Header
	fmt.Printf("%-10s %-18s %-14s\n", "NAME", "IP", "SSH STATUS")
	fmt.Println("────────── ────────────────── ──────────────")

	var serversToExtract []extractableServer

	for name, ip := range servers {
		ping := getStatus(ip)
		var status string

		switch ping {
		case true:
			status = "reachable"
		case false:
			status = "unreachable"
		}

		serversToExtract = append(serversToExtract, extractableServer{Name: name, IP: ip, Status: status})

		fmt.Printf("%-10s %-18s %-14s\n", name, ip, status)
		fmt.Println("────────── ────────────────── ──────────────")
	}

	switch outputType {
	case "json":
		if err := extractToJSON(serversToExtract); err != nil {
			log.Printf("[❌] Failed to extract to JSON: %s\n", err)
			return err
		}
	case "csv":
		// Implement CSV output formatting
	default:
		// Default output already printed above
	}
	return nil
}

func extractToJSON(servers []extractableServer) error {
	filename := "servers_output.json"
	fmt.Printf("[i] Extracting to JSON file: %s\n", filename)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create JSON file: %s", err)
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(servers); err != nil {
		return fmt.Errorf("Failed to encode servers to JSON: %s", err)
	}
	return nil
}
