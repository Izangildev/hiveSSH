package logic

import (
	"bytes"
	"fmt"
	"hivessh/env"

	"github.com/melbahja/goph"
)

func Run(command, identifier string) error {

	exists, kind := serverExists(identifier)
	if !exists {
		fmt.Printf("[❌] Server '%s' not found in database\n", identifier)
	} else {
		fmt.Printf("[✅] Server '%s' found by %s\n", identifier, kind)
	}

	var ip string
	var stdout, stderr bytes.Buffer
	var port int
	var user string

	switch kind {
	case "name":
		ip = servers[identifier].IP
		port = servers[identifier].Port
		user = servers[identifier].User
	case "IP":
		ip = identifier
		for _, server := range servers {
			if server.IP == identifier {
				port = server.Port
				user = server.User
				break
			}
		}

	default:
		return fmt.Errorf("invalid identifier type")
	}

	// Start new ssh connection with private key.
	auth, err := goph.Key(env.Private_key, "")
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	// client, err := goph.New(user, fmt.Sprintf("%s:%d", ip, port), auth)
	client, err := goph.New(user, ip, auth)
	if err != nil {
		return fmt.Errorf("failed to connect to %s:%d: %w", ip, port, err)
		// return fmt.Errorf("failed to connect to %s: %w", ip, err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}

	defer session.Close()

	// session.Stdout stores the address to write output
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		return fmt.Errorf("command failed: %s\n[stderr]: %s", err, stderr.String())
	}

	if stdout.Len() > 0 {
		fmt.Printf("[✅] Output:\n%s\n", stdout.String())
	} else {
		fmt.Println("[✅] Command executed successfully with no output.")
	}

	return nil
}
