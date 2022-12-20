package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/danilomarques/secretumcli/pb"
)

type Shell struct {
	passwordClient pb.PasswordClient
	scanner        *bufio.Scanner
	token          string
}

func NewShell(passwordClient pb.PasswordClient, token string) *Shell {
	scanner := bufio.NewScanner(os.Stdin)
	return &Shell{
		passwordClient: passwordClient,
		scanner:        scanner,
		token:          token,
	}
}

func (s *Shell) Run() {
	for {
		fmt.Print(">> ")
		var cmd string
		if s.scanner.Scan() {
			cmd = s.scanner.Text()
		}
		switch cmd {
		}
	}
}
