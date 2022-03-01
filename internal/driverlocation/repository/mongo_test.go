package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/s3f4/locationmatcher/internal/driverlocation/models"
	"github.com/s3f4/locationmatcher/pkg/log"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

func TestMain(m *testing.M) {

	os.Setenv("DRIVER_LOCATION_DATABASE", "driver_location")
	os.Setenv("DRIVER_LOCATION_COLLECTION", "driver_location")

	// Setup
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	environmentVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=root",
		"MONGO_INITDB_ROOT_PASSWORD=password",
	}

	resource, err := pool.Run("mongo", "5.0", environmentVariables)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		var err error
		db, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		return db.Ping(context.TODO(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// seed data
	log.Info("Tests are starting...")
	DB = "driver_location"
	Collection = "driver_location"

	// Run tests
	exitCode := m.Run()

	// Teardown
	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	log.Info("Tests Finished")
	// Exit
	os.Exit(exitCode)
}

func Test_Mongo(t *testing.T) {
	ctx := context.Background()
	repo, err := NewRepository("mongo", db)
	assert.Nil(t, err)

	driverLocations := []*models.DriverLocation{
		// galata
		{
			Location: models.Location{
				Type:        "Point",
				Coordinates: []interface{}{28.97413088610361, 41.025651081666744},
			},
		},
		// ayasofya
		{
			Location: models.Location{
				Type:        "Point",
				Coordinates: []interface{}{28.979986854317975, 41.00858654897259},
			},
		},
		// iu
		// {
		// 	Location: models.Location{
		// 		Type:        "Point",
		// 		Coordinates: []interface{}{28.9605116156308, 41.01189519061322},
		// 	},
		// },
	}
	// insert driverLocations
	err = repo.UpsertBulk(ctx, driverLocations)
	// assert error is nil
	assert.Nil(t, err)

	err = repo.CreateIndex(ctx, "location", "2dsphere")
	assert.Nil(t, err)

	// it should return ayasofya -> galata
	locations, err := repo.Find1(ctx, &models.Query{
		Location: models.Location{
			Type:        "Point",
			Coordinates: []interface{}{28.9605116156308, 41.01189519061322},
		},
		MinDistance: 0,
		MaxDistance: 10000,
	})

	assert.Nil(t, err)
	assert.Len(t, locations, 2)

	fmt.Printf("%#v\n", *locations[0].MongoDistance)
	fmt.Printf("%#v\n", locations[0].Distance)
	fmt.Printf("%#v\n", *locations[1].MongoDistance)
	fmt.Printf("%#v\n", locations[1].Distance)
	coords, _ := locations[0].Location.Coordinates.(primitive.A)
	assert.Equal(t, coords[0], driverLocations[1].Location.Coordinates.([]interface{})[0])

	// it s
	// assert.Equal(t, locations[0].Location.Coordinates)

	locations, err = repo.Find1(ctx, &models.Query{
		Location: models.Location{
			Type:        "Point",
			Coordinates: []interface{}{28.9605116156308, 41.01189519061322},
		},
		MinDistance: 0,
		MaxDistance: 10000,
	})

	assert.Nil(t, err)
	assert.Len(t, locations, 2)

	err = repo.DropIfExists(ctx)
	assert.Nil(t, err)
}

// func Test_MongoRepository_Find(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	id1 := primitive.NewObjectID()
// 	id2 := primitive.NewObjectID()

// 	point1 := primitive.M{
// 		"0": 40.94001079,
// 		"1": 29.00077262,
// 	}

// 	point2 := primitive.M{
// 		"0": 35.9421,
// 		"1": 30.500,
// 	}

// 	mt.Run("success", func(mt *mtest.T) {
// 		driverLocationRepository := NewMongoRepository(mt.Client)
// 		first := mtest.CreateCursorResponse(1, "mongo.record", mtest.FirstBatch, bson.D{
// 			{"_id", id1},
// 			{"location", bson.D{
// 				{"type", "Point"},
// 				{"coordinates", point1},
// 			}},
// 		})

// 		second := mtest.CreateCursorResponse(1, "mongo.record", mtest.NextBatch, bson.D{
// 			{"_id", id2},
// 			{"location", bson.D{
// 				{"type", "Point"},
// 				{"coordinates", point2},
// 			}},
// 		})
// 		killCursors := mtest.CreateCursorResponse(0, "mongo.record", mtest.NextBatch)
// 		mt.AddMockResponses(first, second, killCursors)

// 		records, err := driverLocationRepository.Find(context.Background(), &models.Query{})

// 		resultRecords := []*models.DriverLocation{
// 			{
// 				ID: id1,
// 				Location: models.Location{
// 					Type:        "Point",
// 					Coordinates: []float64{40.94001079, 29.00077262},
// 				},
// 			},
// 			{
// 				ID: id2,
// 				Location: models.Location{
// 					Type:        "Point",
// 					Coordinates: []float64{35.9421, 30.500},
// 				},
// 			},
// 		}

// 		assert.Nil(t, err)
// 		assert.Equal(t, records, resultRecords)
// 	})

// 	mt.Run("command_fail", func(mt *mtest.T) {
// 		driverLocationRepository := NewMongoRepository(mt.Client)
// 		mt.AddMockResponses(bson.D{{"ok", 0}})
// 		_, err := driverLocationRepository.Find(context.Background(), &models.Query{})

// 		assert.NotNil(t, err)
// 		assert.IsType(t, err, mongo.CommandError{})
// 	})

// 	mt.Run("cursor_fail", func(mt *mtest.T) {
// 		driverLocationRepository := NewMongoRepository(mt.Client)

// 		first := mtest.CreateCursorResponse(0, "mongo.record", mtest.FirstBatch, bson.D{
// 			{"_id", 1},
// 		})
// 		mt.AddMockResponses(first)

// 		records, err := driverLocationRepository.Find(context.Background(), &models.Query{})

// 		assert.Nil(t, records)
// 		assert.NotNil(t, err)
// 	})
// }
