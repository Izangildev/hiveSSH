package logic

import (
	"fmt"
)

func Join(serverName, ip, user, description string, port int) error {
	existsName, _ := serverExists(serverName)
	existsIP, _ := serverExists(ip)

	if existsName || existsIP {
		return fmt.Errorf("server %s or IP %s already stored", serverName, ip)
	}

	serverToStore := ServerInfo{
		IP:          ip,
		User:        user,
		Port:        port,
		Group:       []string{},
		Description: description,
	}

	servers[serverName] = serverToStore
	SaveServers()
	return nil
}
