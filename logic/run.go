package logic

import (
	"bytes"
	"fmt"
	"hivessh/env"
	"os"

	"golang.org/x/crypto/ssh"
)

func Run(command, identifier string) error {
	exists, kind := serverExists(identifier)
	if !exists {
		fmt.Printf("[❌] Server '%s' not found in database\n", identifier)
		return fmt.Errorf("server not found")
	} else {
		fmt.Printf("[✅] Server '%s' found by %s\n", identifier, kind)
	}

	var ip string
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

	key, err := os.ReadFile(env.Private_key)
	if err != nil {
		return fmt.Errorf("unable to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use a proper host key callback
		Timeout:         20 * 1e9,
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
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
