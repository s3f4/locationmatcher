package repository

import (
	"context"

	"github.com/s3f4/locationmatcher/internal/driverlocation/models"
)

// Repository ..
type Repository interface {
	UpsertBulk(context.Context, []*models.DriverLocation) error
	Find(context.Context, *models.Query) ([]*models.DriverLocation, error)
	Find1(context.Context, *models.Query) ([]*models.DriverLocation, error)
	DropIfExists(context.Context) error
	CreateIndex(context.Context, string, string) error
	Migrate(context.Context)
}
