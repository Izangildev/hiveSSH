package group

import (
	"crypto/md5"
	"fmt"
	"hivessh/env"
	"hivessh/logic"
	"time"
)

func Create(groupName string) error {
	if !logic.ExistGroupsFile(env.GroupsFile) {
		return fmt.Errorf("groups file does not exist")
	}

	if _, exists := logic.Groups[groupName]; exists {
		return fmt.Errorf("group '%s' already exists", groupName)
	}

	var id string = createID()

	var group logic.GroupInfo = logic.GroupInfo{
		Id:          id,
		Description: "",
		Members:     []string{},
	}

	logic.Groups[groupName] = group
	logic.SaveGroups()
	return nil
}

// Function that creates a hash id for groups
func createID() string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
