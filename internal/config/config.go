package config

import (
	"os"
	"path/filepath"
)

var (
	PrivateKey  string
	DataDir     string
	ServersFile string
	GroupsFile  string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("hivessh: cannot determine home directory: " + err.Error())
	}

	PrivateKey  = filepath.Join(home, ".ssh", "id_rsa")
	DataDir     = `C:\Users\izang\Desktop\PROYECTOS\hiveSSH`
	ServersFile = filepath.Join(DataDir, "servers.json")
	GroupsFile  = filepath.Join(DataDir, "groups.json")
}

// InitDataDir creates the hivessh data directory if it does not exist.
func InitDataDir() error {
	return os.MkdirAll(DataDir, 0755)
}
