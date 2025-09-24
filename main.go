package main

import (
	"hivessh/cmd"
	_ "hivessh/cmd/group" // <-- Fuerza la inicialización de los comandos group
	"hivessh/env"
	"hivessh/logic"
	"hivessh/logic/group"
)

func main() {
	group.LoadGroups(env.GroupsFile)
	logic.LoadServers(env.ServersFile)
	cmd.Execute()
}
