package group

import (
	"fmt"
	"hivessh/logic"
)

func JoinServer(groupname, serverID string) error {
	serverExists, _ := logic.ServerExists(serverID)
	if !GroupExists(groupname) || !serverExists {
		return fmt.Errorf("group %s or server %s does not exist", groupname, serverID)
	}

	groupID := groups[groupname].Id
	// Check if server is already in the group
	for _, member := range groups[groupname].Members {
		if member == serverID {
			return fmt.Errorf("server %s is already a member of group %s", serverID, groupname)
		}
	}

	group := groups[groupname]
	group.Members = append(group.Members, serverID)
	groups[groupname] = group

	// Falta añadir el grupo al servidor
	var server logic.ServerInfo
	var serverName string
	for name, srv := range logic.Servers {
		if srv.Id == serverID {
			server = srv
			serverName = name
			break
		}
	}
	server.Groups = append(server.Groups, groupID)
	logic.Servers[serverName] = server

	SaveGroups()
	logic.SaveServers()
	return nil
}
