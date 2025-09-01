package logic

import (
	"encoding/json"
	"fmt"
	"hivessh/env"
	"os"
)

var servers = make(map[string]string)

func existServersFile(serversFile string) bool {
	_, err := os.Stat(serversFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println("Error al comprobar el archivo:", err)
		return false
	}
	return true
}

// Guarda servers en el JSON
func SaveServers() {
	data, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		fmt.Println("Error al transformar en JSON:", err)
		return
	}

	err = os.WriteFile(env.ServersFile, data, 0644)
	if err != nil {
		fmt.Println("Error al escribir el fichero:", err)
		return
	}
}

func LoadServers(serversFile string) {
	if !existServersFile(serversFile) {
		return
	}

	data, err := os.ReadFile(serversFile)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	if len(data) == 0 {
		return
	}

	err = json.Unmarshal(data, &servers)
	if err != nil {
		fmt.Println("Error al parsear JSON:", err)
		return
	}
}
