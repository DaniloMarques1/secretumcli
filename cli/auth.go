package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/danilomarques/secretumcli/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// request user information and try register
// the user on the server
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

// request user information and try to sign in
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
		errStatus, ok := status.FromError(err)
		if ok && errStatus.Code() == codes.PermissionDenied {
			if err := a.handlePasswordExpired(email, password); err != nil {
				return "", err
			}

			log.Printf("Your password was updated. Please try sign in again\n")
			os.Exit(1)
		} else {
			return "", err
		}
	}

	return response.AccessToken, nil
}

// it will ask for a new password if the current got expired
func (a *Auth) handlePasswordExpired(email, oldPassword string) error {
	newPassword, err := ReadPassword("Please provide a new password: ")
	if err != nil {
		return err
	}

	if newPassword == oldPassword {
		return errors.New("The new password should be different than your current one")
	}

	req := &pb.UpdateMasterRequest{
		Email:       email,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	_, err = a.client.UpdateMaster(
		context.Background(),
		req,
	)

	if err != nil {
		return err
	}

	return nil
}
