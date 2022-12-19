package cli

import (
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func ReadPassword(label string) (string, error) {
	fmt.Print(label)
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // cleaning input

	return string(b), nil
}
