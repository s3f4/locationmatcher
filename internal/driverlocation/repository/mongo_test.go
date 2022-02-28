package repository

import (
	"context"
	"os"
	"testing"

	"github.com/s3f4/locationmatcher/internal/driverlocation/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMain(m *testing.M) {
	os.Setenv("DRIVER_LOCATION_DATABASE", "test_database")
	os.Setenv("DRIVER_LOCATION_COLLECTION", "test_driver_location")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_MongoRepository_Find(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	point1 := primitive.M{
		"0": 40.94001079,
		"1": 29.00077262,
	}

	point2 := primitive.M{
		"0": 35.9421,
		"1": 30.500,
	}

	mt.Run("success", func(mt *mtest.T) {
		driverLocationRepository := NewMongoRepository(mt.Client)
		first := mtest.CreateCursorResponse(1, "mongo.record", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"location", bson.D{
				{"type", "Point"},
				{"coordinates", point1},
			}},
		})

		second := mtest.CreateCursorResponse(1, "mongo.record", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"location", bson.D{
				{"type", "Point"},
				{"coordinates", point2},
			}},
		})
		killCursors := mtest.CreateCursorResponse(0, "mongo.record", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		records, err := driverLocationRepository.Find(context.Background(), &models.Query{})

		resultRecords := []*models.DriverLocation{
			{
				ID: id1,
				Location: models.Location{
					Type:        "Point",
					Coordinates: []float64{40.94001079, 29.00077262},
				},
			},
			{
				ID: id2,
				Location: models.Location{
					Type:        "Point",
					Coordinates: []float64{35.9421, 30.500},
				},
			},
		}

		assert.Nil(t, err)
		assert.Equal(t, records, resultRecords)
	})

	mt.Run("command_fail", func(mt *mtest.T) {
		driverLocationRepository := NewMongoRepository(mt.Client)
		mt.AddMockResponses(bson.D{{"ok", 0}})
		_, err := driverLocationRepository.Find(context.Background(), &models.Query{})

		assert.NotNil(t, err)
		assert.IsType(t, err, mongo.CommandError{})
	})

	mt.Run("cursor_fail", func(mt *mtest.T) {
		driverLocationRepository := NewMongoRepository(mt.Client)

		first := mtest.CreateCursorResponse(0, "mongo.record", mtest.FirstBatch, bson.D{
			{"_id", 1},
		})
		mt.AddMockResponses(first)

		records, err := driverLocationRepository.Find(context.Background(), &models.Query{})

		assert.Nil(t, records)
		assert.NotNil(t, err)
	})
}
