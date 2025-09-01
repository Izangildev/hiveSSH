package logic

import (
	"fmt"
)

func serverExists(serverName, ipAddress string) bool {
	for name, ip := range servers {
		if serverName == name || ipAddress == ip {
			return true
		}
	}
	return false
}

func Join(serverName, ip string) {
	if serverExists(serverName, ip) {
		fmt.Println("This server is already loaded.")
		return
	}

	servers[serverName] = ip
	fmt.Println("Server", serverName, "saved with IP -->", ip)
	SaveServers()
}
