package main

import (
	"errors"
	"log"
	"os"

	"github.com/danilomarques/secretumcli/cli"
)

func main() {
	arg, err := parseArguments()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := openGrpcConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := cli.NewCli(conn)
	c.Run(arg)
}

func parseArguments() (string, error) {
	args := os.Args
	if len(args) != 2 {
		return "", errors.New("Invalid number of arguments. Use the \"help\" command to see the usage")
	}

	return args[1], nil
}
