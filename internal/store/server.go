package store

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"hivessh/internal/config"
	"log"
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

func createID() string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func getStatus(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":22", 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

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

func SaveServers() {
	data, err := json.MarshalIndent(Servers, "", "  ")
	if err != nil {
		fmt.Printf("[❌] Failed to convert in JSON: %s\n", err)
		return
	}
	err = os.WriteFile(config.ServersFile, data, 0644)
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

func Join(serverName, ip, user, description string, port int) error {
	if existsName, _ := ServerExists(serverName); existsName {
		return fmt.Errorf("a server named '%s' already exists", serverName)
	}
	if existsIP, _ := ServerExists(ip); existsIP {
		return fmt.Errorf("IP '%s' is already registered", ip)
	}
	var id string = createID()
	serverToStore := ServerInfo{
		Id:          serverName + id,
		IP:          ip,
		User:        user,
		Port:        port,
		Groups:      []string{},
		Description: description,
	}
	Servers[serverName] = serverToStore
	SaveServers()
	return nil
}

type extractableServer struct {
	Name   string
	Id     string
	IP     string
	Status string
}

func List(outputType string) error {
	fmt.Printf("%-10s %-18s %-14s %-35s\n", "   NAME  ", "        IP      ", "  SSH STATUS ", "              ID              ")
	fmt.Println("────────── ────────────────── ────────────── ─────────────────────────────────")

	var serversToExtract []extractableServer

	for name := range Servers {
		ip := Servers[name].IP
		id := Servers[name].Id
		var status string
		if getStatus(ip) {
			status = "reachable"
		} else {
			status = "unreachable"
		}
		serversToExtract = append(serversToExtract, extractableServer{Name: name, IP: ip, Status: status, Id: id})
		fmt.Printf("%-10s %-18s %-14s %-35s\n", name, ip, status, id)
		fmt.Println("────────── ────────────────── ────────────── ─────────────────────────────────")
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
	if err := writer.Write([]string{"Name", "Id", "IP", "Status"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %s", err)
	}
	for _, server := range servers {
		record := []string{server.Name, server.Id, server.IP, server.Status}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %s", err)
		}
	}
	return nil
}
