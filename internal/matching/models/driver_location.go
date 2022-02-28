package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DriverLocation holds the driver's location data
type DriverLocation struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Location Location           `json:"location" bson:"location"`
	Distance float64            `json:"distance" bson:"-"`
}
