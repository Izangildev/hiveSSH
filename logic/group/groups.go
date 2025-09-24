package group

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

var groups = make(map[string]GroupInfo)

func existGroupsFile(groupsFile string) bool {
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
	data, err := json.MarshalIndent(groups, "", "  ")
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
	if !existGroupsFile(groupsFile) {
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

	err = json.Unmarshal(data, &groups)
	if err != nil {
		fmt.Printf("[❌] Failed to parse groups JSON: %s\n", err)
		return
	}
}
