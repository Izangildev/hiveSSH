package ssh

import (
	"bytes"
	"fmt"
	"hivessh/internal/config"
	"hivessh/internal/store"
	"os"
	"strings"
	"sync"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

// serverTarget holds the resolved connection details for a server.
type serverTarget struct {
	name string
	ip   string
	port int
	user string
}

// resolveServer finds a server by name, IP or ID and returns its connection details.
func resolveServer(identifier string) (serverTarget, error) {
	exists, kind := store.ServerExists(identifier)
	if !exists {
		return serverTarget{}, fmt.Errorf("server '%s' not found in database", identifier)
	}

	var t serverTarget
	switch kind {
	case "name":
		srv := store.Servers[identifier]
		t = serverTarget{name: identifier, ip: srv.IP, port: srv.Port, user: srv.User}
	case "IP":
		for name, srv := range store.Servers {
			if srv.IP == identifier {
				t = serverTarget{name: name, ip: identifier, port: srv.Port, user: srv.User}
				break
			}
		}
	case "ID":
		for name, srv := range store.Servers {
			if srv.Id == identifier {
				t = serverTarget{name: name, ip: srv.IP, port: srv.Port, user: srv.User}
				break
			}
		}
	default:
		return serverTarget{}, fmt.Errorf("unknown identifier type: %s", kind)
	}
	return t, nil
}

// newSSHConfig reads the private key and builds a gossh.ClientConfig.
func newSSHConfig(user string) (*gossh.ClientConfig, error) {
	key, err := os.ReadFile(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %w", err)
	}
	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}
	return &gossh.ClientConfig{
		User:            user,
		Auth:            []gossh.AuthMethod{gossh.PublicKeys(signer)},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}, nil
}

// Run executes a command on a single server identified by name, IP or ID.
func Run(command, identifier string) error {
	target, err := resolveServer(identifier)
	if err != nil {
		return err
	}
	fmt.Printf("[✅] Server '%s' found\n", target.name)

	cfg, err := newSSHConfig(target.user)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", target.ip, target.port)
	client, err := gossh.Dial("tcp", addr, cfg)
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

// runOnMemberOutput executes a command on one group member and returns formatted output.
func runOnMemberOutput(command, identifier string) string {
	var output strings.Builder

	target, err := resolveServer(identifier)
	if err != nil {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] %s\n", identifier, err))
		return output.String()
	}

	cfg, err := newSSHConfig(target.user)
	if err != nil {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] %s\n", target.name, err))
		return output.String()
	}

	addr := fmt.Sprintf("%s:%d", target.ip, target.port)
	client, err := gossh.Dial("tcp", addr, cfg)
	if err != nil {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] Failed to connect to %s: %s\n", target.name, addr, err))
		return output.String()
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		output.WriteString(fmt.Sprintf("[%s]\n\t[❌] Failed to create SSH session: %s\n", target.name, err))
		return output.String()
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	output.WriteString(fmt.Sprintf("[%s]\n", strings.ToUpper(target.name)))
	if err := session.Run(command); err != nil {
		output.WriteString(fmt.Sprintf("\t[❌] Error executing command\n\t[ERROR DETAILS] %s\n", strings.ReplaceAll(stderr.String(), "\n", "\n\t")))
	} else if stdout.Len() > 0 {
		output.WriteString(fmt.Sprintf("\t[✅] Command: %s executed successfully\n", command))
		output.WriteString("\t[✅] Output:\n")
		output.WriteString(fmt.Sprintf("\t%s\n", strings.ReplaceAll(stdout.String(), "\n", "\n\t")))
	} else {
		output.WriteString(fmt.Sprintf("\t[✅] Command: %s executed successfully with no output\n", command))
	}
	return output.String()
}

// RunGroup executes a command concurrently on all members of a group.
func RunGroup(command, groupname string) error {
	if !store.GroupExists(groupname) {
		return fmt.Errorf("group '%s' does not exist", groupname)
	}
	members := store.Groups[groupname].Members
	if len(members) == 0 {
		return fmt.Errorf("group '%s' has no members", groupname)
	}

	var wg sync.WaitGroup
	outputCh := make(chan string, len(members))

	for _, member := range members {
		wg.Add(1)
		go func(member string) {
			defer wg.Done()
			outputCh <- runOnMemberOutput(command, member)
		}(member)
	}

	wg.Wait()
	close(outputCh)

	for out := range outputCh {
		fmt.Print(out)
	}
	return nil
}
