package server

import (
	"context"
	"fmt"

	"github.com/s3f4/locationmatcher/internal/matching/client"
)

type Server interface {
	Start(context.Context)
}

func NewServer(server string) (Server, error) {
	switch server {
	case "http":
		return &httpServer{
			client: client.NewAPIClient(),
		}, nil

	case "grpc":
		return &grpcServer{}, nil

	default:
		return nil, fmt.Errorf("no such server %s", server)
	}
}
