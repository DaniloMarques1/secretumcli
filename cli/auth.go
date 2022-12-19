package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/danilomarques/secretumcli/pb"
)

var (
	ErrPasswordDoesNotMatch = errors.New("Password does not match")
	ErrCreatingMaster       = errors.New("Something went wrong while creating a master")
)

type Auth struct {
	client  pb.MasterClient
	scanner *bufio.Scanner
}

func NewAuth(client pb.MasterClient) *Auth {
	scanner := bufio.NewScanner(os.Stdin)
	return &Auth{client: client, scanner: scanner}
}

func (a *Auth) SignUp() error {
	fmt.Print("Please provide an email address: ")
	var email string
	if a.scanner.Scan() {
		email = a.scanner.Text()
	}

	password, err := ReadPassword("Please provide a (good) password: ")
	if err != nil {
		return err
	}

	confirmPassword, err := ReadPassword("Type your password one more time please: ")
	if err != nil {
		return err
	}

	if password != confirmPassword {
		return ErrPasswordDoesNotMatch
	}

	response, err := a.client.SaveMaster(
		context.Background(),
		&pb.CreateMasterRequest{Email: email, Password: password},
	)
	if err != nil {
		return err
	}

	if !response.OK {
		return ErrCreatingMaster
	}

	return nil
}

func (a *Auth) SignIn() (string, error) {
	fmt.Print("Please provide an email address: ")
	var email string
	if a.scanner.Scan() {
		email = a.scanner.Text()
	}

	password, err := ReadPassword("Please provide your password: ")
	if err != nil {
		return "", err
	}

	response, err := a.client.AuthenticateMaster(
		context.Background(),
		&pb.AuthMasterRequest{Email: email, Password: password},
	)

	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}
