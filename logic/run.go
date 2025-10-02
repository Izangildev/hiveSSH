package logic

import (
	"bytes"
	"fmt"
	"hivessh/env"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

func Run(command, identifier string) error {
	exists, kind := ServerExists(identifier)
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
		ip = Servers[identifier].IP
		port = Servers[identifier].Port
		user = Servers[identifier].User
	case "IP":
		ip = identifier
		for _, server := range Servers {
			if server.IP == identifier {
				port = server.Port
				user = server.User
				break
			}
		}
	case "ID":
		for _, server := range Servers {
			if server.Id == identifier {
				ip = server.IP
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

func runOnMember(command, identifier string) error {
	exists, kind := ServerExists(identifier)
	if !exists {
		return fmt.Errorf("server '%s' not found", identifier)
	}

	var ip string
	var port int
	var user string
	var name string

	switch kind {
	case "name":
		ip = Servers[identifier].IP
		port = Servers[identifier].Port
		user = Servers[identifier].User
		name = identifier
	case "IP":
		ip = identifier
		for key, server := range Servers {
			if server.IP == identifier {
				port = server.Port
				user = server.User
				name = key
				break
			}
		}
	case "ID":
		for key, server := range Servers {
			if server.Id == identifier {
				ip = server.IP
				port = server.Port
				user = server.User
				name = key
				break
			}
		}
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

	fmt.Printf("[✅] Command: %s executed successfully in %s\n", command, name)

	if stdout.Len() > 0 {
		fmt.Printf("[✅] Output:\n%s\n", stdout.String())
	}

	return nil
}

func RunGroup(command, groupname string) error {
	if !GroupExists(groupname) {
		return fmt.Errorf("group '%s' does not exist", groupname)
	}

	members := Groups[groupname].Members
	if len(members) == 0 {
		return fmt.Errorf("group '%s' has no members", groupname)
	}

	var waitGroup sync.WaitGroup

	for _, member := range members {
		// Capture member name

		waitGroup.Add(1)
		go func(member string) {
			defer waitGroup.Done()
			if err := runOnMember(command, member); err != nil {
				fmt.Printf("[❌] Error executing command on member '%s': %s\n", member, err)
			}
		}(member)
	}

	waitGroup.Wait()
	return nil
}
