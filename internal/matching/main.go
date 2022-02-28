package main

import (
	"context"
	"os"

	"github.com/s3f4/locationmatcher/internal/matching/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := server.NewServer(os.Getenv("SERVER"))
	if err != nil {
		panic(err)
	}

	// Starts server
	server.Start(ctx)
}
