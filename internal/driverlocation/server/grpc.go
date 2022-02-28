package server

import (
	"context"

	"github.com/s3f4/locationmatcher/internal/driverlocation/repository"
	"github.com/s3f4/locationmatcher/pkg/log"
)

type grpcServer struct{}

func (*grpcServer) Start(ctx context.Context, repository repository.Repository) {
	log.Fatal("not implemented")
}
