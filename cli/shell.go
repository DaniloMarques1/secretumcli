package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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
		var input string
		if s.scanner.Scan() {
			input = s.scanner.Text()
		}
		cmd, args, err := s.parseInput(input)
		if err != nil {
			continue
		}

		switch cmd {
		case SAVE:
			if len(args) != 2 {
				fmt.Println("You need to provide a key and a password")
				continue
			}
			key := args[0]
			password := args[1]
			request := &pb.CreatePasswordRequest{
				AccessToken: s.token,
				Key:         key,
				Password:    password,
			}
			response, err := s.passwordClient.SavePassword(context.Background(), request)
			if err != nil || !response.OK {
				fmt.Println("Something went wrong. Please try again")
				continue
			}

			fmt.Println("Password saved successfully")
		case FIND:
			if len(args) != 1 {
				fmt.Println("You need to provide the key you want to look for")
				continue
			}
			key := args[0]
			request := &pb.FindPasswordRequest{
				AccessToken: s.token,
				Key:         key,
			}
			response, err := s.passwordClient.FindPassword(context.Background(), request)
			if err != nil {
				fmt.Println("Something went wrong. Please try again")
				continue
			}

			fmt.Println(response.Password)
		case REMOVE:
			if len(args) != 1 {
				fmt.Println("You need to provide the key you want to remove")
				continue
			}
			key := args[0]
			request := &pb.RemovePasswordRequest{
				AccessToken: s.token,
				Key:         key,
			}
			response, err := s.passwordClient.RemovePassword(context.Background(), request)
			if err != nil || !response.OK {
				fmt.Println("Something went wrong. Please try again")
				continue
			}

			fmt.Println("Password removed successfully")
		case GENERATE:
			if len(args) != 2 {
				fmt.Println("You need to a provide the key and a keyphrase that will be used to generate the random password")
				continue
			}
			key := args[0]
			keyphrase := args[1]
			request := &pb.GeneratePasswordRequest{
				AccessToken: s.token,
				Key:         key,
				Keyphrase:   keyphrase,
			}
			_, err := s.passwordClient.GeneratePassword(context.Background(), request)

			if err != nil {
				fmt.Println("Something went wrong. Please try again")
				continue
			}

			fmt.Println("Password generated successfully")
		case KEYS:
			if len(args) != 0 {
				fmt.Println("Invalid number of arguments")
				continue
			}
			request := &pb.FindKeysRequest{AccessToken: s.token}
			response, err := s.passwordClient.FindKeys(context.Background(), request)
			if err != nil {
				fmt.Println("Something went wrong. Please try again")
				continue
			}

			if len(response.GetKeys()) == 0 {
				fmt.Println("You have not saved any passwords yet")
				continue
			}

			for _, key := range response.GetKeys() {
				fmt.Printf("- %v\n", key)
			}
		case EXIT:
			os.Exit(1)
		case HELP:
			s.Usage()
		default:
			fmt.Println("Command not found")
		}
	}
}

func (s *Shell) parseInput(input string) (string, []string, error) {
	slice := strings.Split(input, " ")
	if len(slice) == 0 {
		return "", nil, errors.New("Wrong input provided")
	}
	cmd := slice[0]
	var args []string
	if len(slice) > 1 {
		args = slice[1:]
	}

	return cmd, args, nil
}

func (s *Shell) Usage() {
	fmt.Println("save      - save a new password example: save passwordkey password")
	fmt.Println("remove    - removes a password example: remove passwordkey")
	fmt.Println("keys      - list all the keys of saved passwords example: keys")
	fmt.Println("update    - update a password example: update passwordkey newpassword")
	fmt.Println("find      - finds the password associated with the given key example: find passwordkey")
	fmt.Println("generate  - will generate a random password example: generate passwordkey passwordkeyphrase")
	fmt.Println("exit      - finish the program")
}
