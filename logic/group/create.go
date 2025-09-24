package group

import (
	"crypto/md5"
	"fmt"
	"hivessh/env"
	"time"
)

func Create(groupName string) error {
	if !existGroupsFile(env.GroupsFile) {
		return fmt.Errorf("groups file does not exist")
	}

	if _, exists := groups[groupName]; exists {
		return fmt.Errorf("group '%s' already exists", groupName)
	}

	var id string = createID()

	var group GroupInfo = GroupInfo{
		Id:          id,
		Description: "",
		Members:     []string{},
	}

	groups[groupName] = group
	SaveGroups()
	return nil
}

// Function that creates a hash id for groups
func createID() string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
