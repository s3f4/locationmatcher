package repository

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/s3f4/locationmatcher/internal/driverlocation/models"
	"github.com/s3f4/locationmatcher/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DB ...
	DB = os.Getenv("DRIVER_LOCATION_DATABASE")
	// location collection
	locationCollection = os.Getenv("DRIVER_LOCATION_COLLECTION")
)

type mongoRepository struct {
	client *mongo.Client
}

// NewMongoRepository returns an DriverLocationRepository object
func NewMongoRepository(client *mongo.Client) Repository {
	return &mongoRepository{
		client,
	}
}

func (r *mongoRepository) getCollection() *mongo.Collection {
	return r.client.Database(DB).Collection(locationCollection)
}

func (r *mongoRepository) Find(ctx context.Context, query *models.Query) ([]*models.DriverLocation, error) {
	collection := r.getCollection()

	coords, _ := query.Location.Coordinates.([]interface{})
	filter := bson.D{
		{
			Key: "location",
			Value: bson.D{
				{
					Key: "$nearSphere", Value: bson.D{
						{
							Key:   "$geometry",
							Value: query.Location,
						},
						{Key: "$minDistance", Value: query.MinDistance},
						{Key: "$maxDistance", Value: query.MaxDistance},
					},
				},
			},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	driverLocations := []*models.DriverLocation{}
	for cursor.Next(ctx) {
		var driverLocation models.DriverLocation
		if err := cursor.Decode(&driverLocation); err != nil {
			log.Error("Could not decode driver location")
			return nil, err
		}

		driverLocation.Distance, err = driverLocation.CalculateDistance(
			coords[0].(float64),
			coords[1].(float64),
		)
		if err != nil {
			log.Error("Could not calculate the distance")
			return nil, err
		}
		driverLocations = append(driverLocations, &driverLocation)
	}

	return driverLocations, nil
}

func (r *mongoRepository) Find1(ctx context.Context, query *models.Query) ([]*models.DriverLocation, error) {
	collection := r.getCollection()
	coords, _ := query.Location.Coordinates.([]interface{})
	pipeline := []bson.M{
		{
			"$geoNear": bson.M{
				"key":           "location",
				"distanceField": "mongo_distance",
				"maxDistance":   query.MaxDistance,
				"spherical":     true,
				"near":          query.Location,
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	driverLocations := []*models.DriverLocation{}
	for cursor.Next(ctx) {
		var driverLocation models.DriverLocation
		if err := cursor.Decode(&driverLocation); err != nil {
			log.Error("Could not decode driver location")
			return nil, err
		}

		mDistance := *driverLocation.MongoDistance / 1000
		driverLocation.MongoDistance = &mDistance

		driverLocation.Distance, err = driverLocation.CalculateDistance(
			coords[0].(float64),
			coords[1].(float64),
		)
		if err != nil {
			log.Error("Could not calculate the distance")
			return nil, err
		}

		driverLocations = append(driverLocations, &driverLocation)
	}

	return driverLocations, nil
}

func (r *mongoRepository) UpsertBulk(ctx context.Context, driverLocations []*models.DriverLocation) error {
	collection := r.getCollection()

	models := []mongo.WriteModel{}
	for _, driverLocation := range driverLocations {
		document := bson.D{
			{
				Key: "location", Value: bson.D{
					{Key: "type", Value: driverLocation.Location.Type},
					{Key: "coordinates", Value: driverLocation.Location.Coordinates},
				},
			},
		}

		if !driverLocation.ID.IsZero() {
			models = append(models,
				mongo.NewUpdateOneModel().
					SetFilter(bson.M{"_id": driverLocation.ID}).
					SetUpdate(bson.D{
						{Key: "$set", Value: document},
					}),
			)
		} else {
			models = append(models, mongo.NewInsertOneModel().SetDocument(document))
		}
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := collection.BulkWrite(ctx, models, opts)
	log.Infof("%#v\n", res)
	log.Info(err)
	return err
}

func (r *mongoRepository) DropIfExists(ctx context.Context) error {
	dbNames, err := r.client.ListDatabaseNames(ctx, bson.D{{Key: "name", Value: DB}})
	if err != nil {
		log.Error(err)
		return err
	}

	for _, name := range dbNames {
		if name == DB {
			return r.client.Database(DB).Drop(ctx)
		}
	}

	return nil
}

// CreateIndex creates index for location query
func (r *mongoRepository) CreateIndex(ctx context.Context, key, value string) error {
	collection := r.getCollection()
	model := mongo.IndexModel{
		Keys: bson.M{
			key: value,
		},
	}

	if _, err := collection.Indexes().CreateOne(ctx, model); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (r *mongoRepository) Migrate(ctx context.Context) {
	if err := r.DropIfExists(ctx); err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("./source/coordinates.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ','
	rows, _ := reader.ReadAll()

	driverLocations := []*models.DriverLocation{}
	for i, row := range rows {
		if i == 0 {
			continue
		}

		driverLocation := &models.DriverLocation{}
		latitude, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		longitude, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		driverLocation.Location.Type = "Point"
		driverLocation.Location.Coordinates = [2]float64{latitude, longitude}
		driverLocations = append(driverLocations, driverLocation)
	}

	log.Info(driverLocations)
	if err := r.UpsertBulk(ctx, driverLocations); err != nil {
		log.Fatal(err)
	}

	if err := r.CreateIndex(ctx, "location", "2dsphere"); err != nil {
		log.Fatal(err)
	}
}
