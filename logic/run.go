package logic

import (
	"fmt"
	"hivessh/env"

	"github.com/melbahja/goph"
)

func Run(command, identifier string) error {

	exists, kind := serverExists(identifier)
	if !exists {
		fmt.Printf("[✖] Server '%s' not found in database\n", identifier)
	} else {
		fmt.Printf("[✔] Server '%s' found by %s\n", identifier, kind)
	}

	var ip string

	switch kind {
	case "name":
		ip = servers[identifier]
	case "IP":
		ip = identifier
	default:
		return fmt.Errorf("invalid identifier type")
	}

	fmt.Println("Executing command:", command)

	// Start new ssh connection with private key.
	auth, err := goph.Key(env.Private_key, "")
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	client, err := goph.New("root", ip, auth)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", ip, err)
	}

	// Defer closing the network connection.
	defer client.Close()

	// Execute your command.
	fmt.Printf("[→] Executing command: %s\n", command)
	out, err := client.Run(command)

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	// Get your output as []byte.
	fmt.Println(string(out))
	return nil
}
