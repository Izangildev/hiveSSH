package main

import (
	"hivessh/cmd"
	"hivessh/env"
	"hivessh/logic"
)

func main() {
	logic.LoadServers(env.ServersFile)
	cmd.Execute()
}
