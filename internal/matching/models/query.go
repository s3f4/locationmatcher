package models

import (
	"fmt"
)

var ErrInvalidCoordinates = fmt.Errorf("invalid coordinates")

type Query struct {
	Location    Location `json:"location"`
	MinDistance int64    `json:"minDistance"`
	MaxDistance int64    `json:"maxDistance"`
}

func (q Query) Validate() error {
	if q.Location.Type != "Point" {
		return fmt.Errorf("you must provide a valid GeoJSON type")
	}

	// coordinates must be []interface{}
	coords, ok := q.Location.Coordinates.([]interface{})
	if !ok {
		return ErrInvalidCoordinates
	}

	if len(coords) != 2 {
		return ErrInvalidCoordinates
	}

	longitude, ok := coords[0].(float64)
	if !ok {
		return ErrInvalidCoordinates
	}

	latitude, ok := coords[1].(float64)
	if !ok {
		return ErrInvalidCoordinates
	}

	if longitude > 180 || longitude < -180 {
		return fmt.Errorf("you must provide a valid longitude")
	}

	if latitude > 90 || latitude < -90 {
		return fmt.Errorf("you must provide a valid latitude")
	}

	if q.MinDistance >= q.MaxDistance {
		return fmt.Errorf("maxDistance must be greater then 0 and minDistance")
	}

	return nil
}
