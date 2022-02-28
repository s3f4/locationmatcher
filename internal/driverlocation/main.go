package main

import (
	"context"
	"os"

	"github.com/s3f4/locationmatcher/internal/driverlocation/repository"
	"github.com/s3f4/locationmatcher/internal/driverlocation/server"
	"github.com/s3f4/locationmatcher/pkg/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connClientMap := repository.InitConnecions()

	repoType := os.Getenv("REPOSITORY")
	repo, err := repository.NewRepository(repoType, connClientMap[repoType])
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("MIGRATE") == "true" {
		repo.Migrate(ctx)
	}

	server, err := server.NewServer(os.Getenv("SERVER"))
	if err != nil {
		log.Fatal(err)
	}

	// Starts server
	server.Start(ctx, repo)
}
