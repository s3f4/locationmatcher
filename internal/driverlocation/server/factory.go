package server

import (
	"context"
	"fmt"

	"github.com/s3f4/locationmatcher/internal/driverlocation/repository"
)

type Server interface {
	Start(context.Context, repository.Repository)
}

func NewServer(server string) (Server, error) {
	switch server {
	case "http":
		return &httpServer{}, nil

	case "grpc":
		return &grpcServer{}, nil

	default:
		return nil, fmt.Errorf("no such server %s", server)
	}
}
