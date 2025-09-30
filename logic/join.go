package logic

import (
	"fmt"
)

func Join(serverName, ip, user, description string, port int) error {
	existsName, _ := ServerExists(serverName)
	existsIP, _ := ServerExists(ip)

	if existsName || existsIP {
		return fmt.Errorf("server %s or IP %s already stored", serverName, ip)
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
