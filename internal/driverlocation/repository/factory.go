package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrClientType     = fmt.Errorf("client type error")
	ErrNotImplemented = fmt.Errorf("not implemented")
)

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
