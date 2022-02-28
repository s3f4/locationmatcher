package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrClientType     = fmt.Errorf("client type error")
	ErrNotImplemented = fmt.Errorf("not implemented")
)

// NewRepository is a repository factory that creates a repository
// object with given parameters.
func NewRepository(repository string, client interface{}) (Repository, error) {
	switch repository {
	case mongoKey:
		mongoClient, ok := client.(*mongo.Client)
		if !ok {
			return nil, ErrClientType
		}
		return &mongoRepository{client: mongoClient}, nil

	case elasticKey:
		return &elasticRepository{}, nil

	default:
		return nil, fmt.Errorf("no such repository %s", repository)
	}
}
