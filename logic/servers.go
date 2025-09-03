package logic

import (
	"encoding/json"
	"fmt"
	"hivessh/env"
	"os"
)

var servers = make(map[string]string)

func serverExists(identifier string) (bool, string) {
	for name, ip := range servers {
		switch {
		case identifier == name:
			return true, "name"
		case identifier == ip:
			return true, "IP"
		}
	}
	return false, ""
}

func existServersFile(serversFile string) bool {
	_, err := os.Stat(serversFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("[✖] Failed to find servers file: %s\n", err)

		return false
	}
	return true
}

// Guarda servers en el JSON
func SaveServers() {
	data, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		fmt.Printf("[✖] Failed to convert in JSON: %s\n", err)
		return
	}

	err = os.WriteFile(env.ServersFile, data, 0644)
	if err != nil {
		fmt.Printf("[✖] Failed to write servers file: %s\n", err)
		return
	}
}

func LoadServers(serversFile string) {
	if !existServersFile(serversFile) {
		return
	}

	data, err := os.ReadFile(serversFile)
	if err != nil {
		fmt.Printf("[✖] Failed to read servers file: %s\n", err)
		return
	}

	if len(data) == 0 {
		return
	}

	err = json.Unmarshal(data, &servers)
	if err != nil {
		fmt.Printf("[✖] Failed to parse servers JSON: %s\n", err)
		return
	}
}
