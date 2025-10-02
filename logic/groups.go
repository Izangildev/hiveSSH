package logic

import (
	"encoding/json"
	"fmt"
	"hivessh/env"
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

	err = os.WriteFile(env.GroupsFile, data, 0644)
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
