package logic

import (
	"fmt"
)

func Join(serverName, ip string) {
	existsName, _ := serverExists(serverName)
	existsIP, _ := serverExists(ip)

	if existsName || existsIP {
		fmt.Println("This server is already stored.")
		return
	}

	servers[serverName] = ip
	fmt.Println("Server", serverName, "saved with IP -->", ip)
	SaveServers()
}
