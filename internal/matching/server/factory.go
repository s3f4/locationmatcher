package server

import (
	"context"
	"fmt"
)

type Server interface {
	Start(context.Context)
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
