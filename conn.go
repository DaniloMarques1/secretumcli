package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const SERVER_ADDR = "localhost:8080"

func openGrpcConnection() (*grpc.ClientConn, error) {
	// TODO: insecure?
	conn, err := grpc.Dial(SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
