package repository

import (
	"context"
	"os"

	"github.com/s3f4/locationmatcher/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoKey   = "mongo"
	elasticKey = "elastic"
)

// mongoClient is used to connect mongodb, connection will be done one time.
var mongoClient *mongo.Client

// connectMongo connects mongodb with mongodb dsn string
func connectMongo(dsn string) (*mongo.Client, error) {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// InitConnecions starts all connections
func InitConnecions() map[string]interface{} {
	clientMap := map[string]interface{}{}

	var err error
	mongoClient, err = connectMongo(os.Getenv("MONGO_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	clientMap[mongoKey] = mongoClient
	return clientMap
}
