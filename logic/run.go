package logic

import (
	"bytes"
	"fmt"
	"hivessh/env"
	"os"
	"strings"
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

// Output OK type:
// [MASTER]
// [✅] <command> executed successfully
// [✅] Output:
// <output>

// Output ERROR type:
// [MASTER]
// [❌] Error executing command on member '<member_name>'
// [ERROR DETAILS]
// <error details>

func runOnMemberOutput(command, identifier string) string {
	var output strings.Builder

	exists, kind := ServerExists(identifier)
	var name string
	if !exists {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] Error executing command\n\t[ERROR DETAILS] server '%s' not found\n", identifier, identifier))
		return output.String()
	}

	var ip string
	var port int
	var user string

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
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] Error executing command\n\t[ERROR DETAILS] unable to read private key: %s\n", name, err))
		return output.String()
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] Error executing command\n\t[ERROR DETAILS] unable to parse private key: %s\n", name, err))
		return output.String()
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * 1e9,
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] Error executing command\n\t[ERROR DETAILS] failed to connect to %s: %s\n", name, addr, err))
		return output.String()
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] Error executing command\n\t[ERROR DETAILS] failed to create SSH session: %s\n", name, err))
		return output.String()
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	output.WriteString(fmt.Sprintf("[%s]\n", strings.ToUpper(name)))

	if err := session.Run(command); err != nil {
		output.WriteString(fmt.Sprintf("\t[❌] Error executing command\n\t[ERROR DETAILS] %s\n", strings.ReplaceAll(stderr.String(), "\n", "\n\t")))
	} else {
		if stdout.Len() > 0 {
			output.WriteString(fmt.Sprintf("\t[✅] Command: %s executed successfully\n", command))
			output.WriteString("\t[✅] Output:\n")
			output.WriteString(fmt.Sprintf("\t%s\n", strings.ReplaceAll(stdout.String(), "\n", "\n\t")))
			return output.String()
		}
		output.WriteString(fmt.Sprintf("\t[✅] Command: %s executed successfully with no output\n", command))
	}

	return output.String()
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
	// Create a channel to collect output from goroutines
	outputCh := make(chan string, len(members))

	for _, member := range members {
		waitGroup.Add(1)
		go func(member string) {
			defer waitGroup.Done()
			// Run the command on the member and send the output to the channel
			output := runOnMemberOutput(command, member)
			outputCh <- output
		}(member)
	}

	waitGroup.Wait()
	close(outputCh)

	// Print all outputs collected from the channel
	for out := range outputCh {
		fmt.Print(out)
	}

	return nil
}
