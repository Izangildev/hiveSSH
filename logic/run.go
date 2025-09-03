package logic

import (
	"fmt"
	"hivessh/env"
	"log"

	"github.com/melbahja/goph"
)

func Run(command, server string) {

	exists, kind := serverExists(server)
	if !exists {
		fmt.Println("This server is not in DB")
	} else {
		fmt.Println("Found server by", kind)
	}

	fmt.Println("Executing command:", command)

	// Start new ssh connection with private key.
	auth, err := goph.Key(env.Private_key, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := goph.New("root", server, auth)
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
