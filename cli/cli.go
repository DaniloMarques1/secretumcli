package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/danilomarques/secretumcli/pb"
)

type Cli struct {
	passwordClient pb.PasswordClient
	masterClient   pb.MasterClient
	scanner        *bufio.Scanner
	token          string
}

func NewCli(passwordClient pb.PasswordClient, masterClient pb.MasterClient, token string) *Cli {
	scanner := bufio.NewScanner(os.Stdin)
	return &Cli{
		passwordClient: passwordClient,
		masterClient:   masterClient,
		scanner:        scanner,
		token:          token,
	}
}

func (c *Cli) Shell() {
	for {
		fmt.Print(">> ")
		var cmd string
		if c.scanner.Scan() {
			cmd = c.scanner.Text()
		}

		switch cmd {
		case ACCESS:
			auth := NewAuth(c.masterClient)
			token, err := auth.SignIn()
			if err != nil {
				log.Printf("ERR: %v\n", err)
				continue
			}

			c.token = token
		case EXIT:
			os.Exit(1)
		default:
			fmt.Println("Command not found")
		}

	}
}
