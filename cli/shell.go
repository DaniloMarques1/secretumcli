package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/danilomarques/secretumcli/pb"
)

type Shell struct {
	passwordClient pb.PasswordClient
	reader         *bufio.Reader
	token          string
}

func NewShell(passwordClient pb.PasswordClient, token string) *Shell {
	reader := bufio.NewReader(os.Stdin)
	return &Shell{
		passwordClient: passwordClient,
		reader:         reader,
		token:          token,
	}
}

// stats a shell where the user can run commands
func (s *Shell) Run() {
	//var input string
	for {
		fmt.Print(">> ")
		input, err := s.reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println()
				os.Exit(1)
			}
			continue
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
		case CLEAR:
			operatingSystem := runtime.GOOS
			switch operatingSystem {
			case "windows":
				exec.Command("cls").Run()
			default:
				// LINUX or MAC
				cmd := exec.Command("clear")
				out, err := cmd.Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(out))
			}

		default:
			fmt.Println("Command not found")
		}
	}
}

// parse the user input returning the commands and its arguments
func (s *Shell) parseInput(input string) (string, []string, error) {
	input = strings.ReplaceAll(input, "\n", "")
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

// shows availables shell commands
func (s *Shell) Usage() {
	fmt.Println("save      - save a new password example: save passwordkey password")
	fmt.Println("remove    - removes a password example: remove passwordkey")
	fmt.Println("keys      - list all the keys of saved passwords example: keys")
	fmt.Println("update    - update a password example: update passwordkey newpassword")
	fmt.Println("find      - finds the password associated with the given key example: find passwordkey")
	fmt.Println("generate  - will generate a random password example: generate passwordkey passwordkeyphrase")
	fmt.Println("clear     - clears the shell")
	fmt.Println("exit      - finish the program")
}
