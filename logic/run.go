package logic

import (
	"fmt"
	"hivessh/env"
	"log"

	"github.com/melbahja/goph"
)

func Run(command, identifier string) {

	exists, kind := serverExists(identifier)
	if !exists {
		fmt.Println("This server is not in DB")
	} else {
		fmt.Println("Found server by", kind)
	}

	var ip string

	switch kind {
	case "name":
		ip = servers[identifier]
	case "IP":
		ip = identifier
	default:
		log.Fatal("Invalid identifier type")
	}

	fmt.Println("Executing command:", command)

	// Start new ssh connection with private key.
	auth, err := goph.Key(env.Private_key, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := goph.New("root", ip, auth)
	if err != nil {
		log.Fatal(err)
	}

	// Defer closing the network connection.
	defer client.Close()

	// Execute your command.
	out, err := client.Run(command)

	if err != nil {
		log.Fatal(err)
	}

	// Get your output as []byte.
	fmt.Println(string(out))
}
