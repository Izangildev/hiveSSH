package main

import (
	"fmt"
	"hivessh/cmd"
	_ "hivessh/cmd/group" // fuerza la inicialización de los subcomandos group
	"hivessh/internal/config"
	"hivessh/internal/store"
	"os"
)

func main() {
	if err := config.InitDataDir(); err != nil {
		fmt.Fprintf(os.Stderr, "[❌] Failed to create config directory: %s\n", err)
		os.Exit(1)
	}
	store.LoadGroups(config.GroupsFile)
	store.LoadServers(config.ServersFile)
	cmd.Execute()
}
