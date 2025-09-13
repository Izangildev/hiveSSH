package logic

import (
	"encoding/csv"
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
	fmt.Printf("%-10s %-18s %-14s\n", "   NAME  ", "        IP      ", "  SSH STATUS ")
	fmt.Println("────────── ────────────────── ──────────────")

	var serversToExtract []extractableServer

	for name := range servers {
		ip := servers[name].IP
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
		if err := extractToCSV(serversToExtract); err != nil {
			log.Printf("[❌] Failed to extract to CSV: %s\n", err)
			return err
		}
	}
	return nil
}

func extractToJSON(servers []extractableServer) error {
	filename := "servers_output.json"
	fmt.Printf("[i] Extracting to JSON file: %s\n", filename)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %s", err)
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(servers); err != nil {
		return fmt.Errorf("failed to encode servers to JSON: %s", err)
	}
	return nil
}

func extractToCSV(servers []extractableServer) error {
	filename := "servers_output.csv"
	fmt.Printf("[i] Extracting to CSV file: %s\n", filename)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %s", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Name", "IP", "Status"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %s", err)
	}

	// Write server da ta
	for _, server := range servers {
		record := []string{server.Name, server.IP, server.Status}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %s", err)
		}
	}
	return nil
}
