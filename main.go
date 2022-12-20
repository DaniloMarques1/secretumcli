package main

import (
	"errors"
	"log"
	"os"

	"github.com/danilomarques/secretumcli/cli"
	"github.com/danilomarques/secretumcli/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	arg, err := parseArguments()
	if err != nil {
		log.Fatal(err)
	}

	addr := os.Getenv("SECRETUM_SERVER")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	masterClient := pb.NewMasterClient(conn)
	passwordClient := pb.NewPasswordClient(conn)
	c := cli.NewCli(passwordClient, masterClient)
	c.Shell(arg)
}

func parseArguments() (string, error) {
	args := os.Args
	if len(args) != 2 {
		return "", errors.New("Invalid number of arguments")
	}

	return args[1], nil
}
