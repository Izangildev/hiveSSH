package store

import (
	"encoding/json"
	"fmt"
	"hivessh/internal/config"
	"os"
)

type GroupInfo struct {
	Id          string
	Description string
	Members     []string
}

var Groups = make(map[string]GroupInfo)

func ExistGroupsFile(groupsFile string) bool {
	_, err := os.Stat(groupsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("[❌] Failed to find groups file: %s\n", err)
		return false
	}
	return true
}

func SaveGroups() {
	data, err := json.MarshalIndent(Groups, "", "  ")
	if err != nil {
		fmt.Printf("[❌] Failed to save groups: %s\n", err)
		return
	}
	err = os.WriteFile(config.GroupsFile, data, 0644)
	if err != nil {
		fmt.Printf("[❌] Failed to write groups file: %s\n", err)
	}
}

func LoadGroups(groupsFile string) {
	if !ExistGroupsFile(groupsFile) {
		return
	}
	data, err := os.ReadFile(groupsFile)
	if err != nil {
		fmt.Printf("[❌] Failed to read groups file: %s\n", err)
		return
	}
	if len(data) == 0 {
		return
	}
	err = json.Unmarshal(data, &Groups)
	if err != nil {
		fmt.Printf("[❌] Failed to parse groups JSON: %s\n", err)
		return
	}
}

func GroupExists(groupname string) bool {
	_, exists := Groups[groupname]
	return exists
}

func Create(groupName string) error {
	if _, exists := Groups[groupName]; exists {
		return fmt.Errorf("group '%s' already exists", groupName)
	}
	group := GroupInfo{
		Id:          createID(),
		Description: "",
		Members:     []string{},
	}
	Groups[groupName] = group
	SaveGroups()
	return nil
}

func JoinServer(groupname, serverID string) error {
	serverExists, _ := ServerExists(serverID)
	if !GroupExists(groupname) || !serverExists {
		return fmt.Errorf("group %s or server %s does not exist", groupname, serverID)
	}
	groupID := Groups[groupname].Id
	for _, member := range Groups[groupname].Members {
		if member == serverID {
			return fmt.Errorf("server %s is already a member of group %s", serverID, groupname)
		}
	}
	group := Groups[groupname]
	group.Members = append(group.Members, serverID)
	Groups[groupname] = group

	var server ServerInfo
	var serverName string
	for name, srv := range Servers {
		if srv.Id == serverID {
			server = srv
			serverName = name
			break
		}
	}
	server.Groups = append(server.Groups, groupID)
	Servers[serverName] = server

	SaveGroups()
	SaveServers()
	return nil
}
