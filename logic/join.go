package logic

import (
	"fmt"
)

func Join(serverName, ip string) error {
	existsName, _ := serverExists(serverName)
	existsIP, _ := serverExists(ip)

	if existsName || existsIP {
		return fmt.Errorf("server %s or IP %s already stored", serverName, ip)
	}

	servers[serverName] = ip
	fmt.Println("Server", serverName, "saved with IP-->", ip)
	SaveServers()
	return nil
}
