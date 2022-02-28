package repository

import (
	"context"
	"log"

	"github.com/s3f4/locationmatcher/internal/driverlocation/models"
)

type elasticRepository struct{}

func (r *elasticRepository) UpsertBulk(context.Context, []*models.DriverLocation) error {
	return ErrNotImplemented
}

func (r *elasticRepository) Find(context.Context, *models.Query) ([]*models.DriverLocation, error) {
	return nil, ErrNotImplemented
}

func (r *elasticRepository) Find1(context.Context, *models.Query) ([]*models.DriverLocation, error) {
	return nil, ErrNotImplemented
}

func (r *elasticRepository) DropIfExists(context.Context) error {
	return ErrNotImplemented
}

func (r *elasticRepository) CreateIndex(context.Context, string, string) error {
	return ErrNotImplemented
}

func (r *elasticRepository) Migrate(context.Context) {
	log.Fatal(ErrNotImplemented)
}
