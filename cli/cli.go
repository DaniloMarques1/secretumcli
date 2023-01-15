package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/danilomarques/secretumcli/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Cli struct {
	passwordClient pb.PasswordClient
	masterClient   pb.MasterClient
	scanner        *bufio.Scanner
}

func NewCli(conn grpc.ClientConnInterface) *Cli {
	masterClient := pb.NewMasterClient(conn)
	passwordClient := pb.NewPasswordClient(conn)

	scanner := bufio.NewScanner(os.Stdin)
	return &Cli{
		passwordClient: passwordClient,
		masterClient:   masterClient,
		scanner:        scanner,
	}
}

func (c *Cli) Run(arg string) {
	auth := NewAuth(c.masterClient)

	switch arg {
	case ACCESS:
		token, err := auth.SignIn()
		if err != nil {
			errStatus, ok := status.FromError(err)
			if ok {
				log.Fatal(errStatus.Message())
			} else {
				log.Fatalf("ERR: %v\n", err)
			}
		}

		shell := NewShell(c.passwordClient, token)
		shell.Run()
	case REGISTER:
		if err := auth.SignUp(); err != nil {
			errStatus, ok := status.FromError(err)
			if ok {
				log.Fatal(errStatus.Message())
			} else {
				log.Fatal(err)
			}
		}
	case HELP:
		c.Usage()
	default:
		log.Fatalf("Invalid command use the help command to check the usage\n")
	}
}

func (c *Cli) Usage() {
	fmt.Println("access     - will request your e-mail and password")
	fmt.Println("register   - you will create a new master")
}
