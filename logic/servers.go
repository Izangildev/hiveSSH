package logic

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hivessh/env"
	"net"
	"os"
	"time"
)

type ServerInfo struct {
	Id          string
	IP          string
	User        string
	Port        int
	Groups      []string
	Description string
}

var Servers = make(map[string]ServerInfo)

// Function that creates a hash id for groups
func createID() string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// This function checks if the server is reachable via SSH (port 22)
func getStatus(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":22", 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// This function returns if the server exists based on an identifier
// If exists returns the identifier kind
func ServerExists(identifier string) (bool, string) {
	for name, server := range Servers {
		switch {
		case identifier == name:
			return true, "name"
		case identifier == server.IP:
			return true, "IP"
		case identifier == server.Id:
			return true, "ID"
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
		fmt.Printf("[❌] Failed to find servers file: %s\n", err)

		return false
	}
	return true
}

// Save servers to JSON
func SaveServers() {
	data, err := json.MarshalIndent(Servers, "", "  ")
	if err != nil {
		fmt.Printf("[❌] Failed to convert in JSON: %s\n", err)
		return
	}

	err = os.WriteFile(env.ServersFile, data, 0644)
	if err != nil {
		fmt.Printf("[❌] Failed to write servers file: %s\n", err)
		return
	}
}

func LoadServers(serversFile string) {
	if !existServersFile(serversFile) {
		return
	}

	data, err := os.ReadFile(serversFile)
	if err != nil {
		fmt.Printf("[❌] Failed to read servers file: %s\n", err)
		return
	}

	if len(data) == 0 {
		return
	}

	err = json.Unmarshal(data, &Servers)
	if err != nil {
		fmt.Printf("[❌] Failed to parse servers JSON: %s\n", err)
		return
	}
}
