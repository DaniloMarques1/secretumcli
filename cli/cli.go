package cli

import (
	"bufio"
	"log"
	"os"

	"github.com/danilomarques/secretumcli/pb"
)

type Cli struct {
	passwordClient pb.PasswordClient
	masterClient   pb.MasterClient
	scanner        *bufio.Scanner
}

func NewCli(passwordClient pb.PasswordClient, masterClient pb.MasterClient) *Cli {
	scanner := bufio.NewScanner(os.Stdin)
	return &Cli{
		passwordClient: passwordClient,
		masterClient:   masterClient,
		scanner:        scanner,
	}
}

func (c *Cli) Shell(arg string) {
	auth := NewAuth(c.masterClient)

	switch arg {
	case ACCESS:
		token, err := auth.SignIn()
		if err != nil {
			log.Printf("ERR: %v\n", err)
			break
		}

		shell := NewShell(c.passwordClient, token)
		shell.Run()
	case REGISTER:
		if err := auth.SignUp(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Invalid command\n")
	}
}
