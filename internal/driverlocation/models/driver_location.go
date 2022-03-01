package models

import (
	"fmt"
	"math"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidCoordinates = fmt.Errorf("invalid coordinates")

// DriverLocation holds the driver's location data
type DriverLocation struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Location      Location           `json:"location" bson:"location"`
	Distance      float64            `json:"distance" bson:"-"`
	MongoDistance *float64           `json:"mongo_distance,omitempty" bson:"mongo_distance,omitempty"`
}

func toRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// r is the mean radius of Earth in kilometers
const r = 6371

// CalculateDistance calculates the distance of two points with Haversine formula and returns
// distance values in kilometers
// d = r * archav(h)= 2 * r * arcsin(sqrt{h})
// d = 2 * r * arcsin(sqrt{ sin2((lat1-lat2)/2 + cos(lat1)*cos(lat2)*sin2((lng1-lng2)/2) })
func (driverLocation *DriverLocation) CalculateDistance(latitude, longitude float64) (float64, error) {
	coordinates, err := driverLocation.getLocation()
	if err != nil {
		return 0, fmt.Errorf("driver location coordinate error")
	}

	dlLat := coordinates[1]
	dlLng := coordinates[0]

	dLat := toRad(latitude - dlLat)
	dLng := toRad(longitude - dlLng)
	lat1 := toRad(dlLat)
	lat2 := toRad(latitude)

	h := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLng/2), 2)
	d := 2 * r * math.Asin(math.Sqrt(h))

	return d, nil
}

// getLocation converts the coordinates of the driver location appropriate values.
func (driverLocation *DriverLocation) getLocation() ([]float64, error) {
	coords := make([]float64, 2)

	switch coordinates := driverLocation.Location.Coordinates.(type) {
	// convert to float64 array if coordinates comes from mongodb
	case primitive.A:
		for i, coord := range coordinates {
			var ok bool
			if coords[i], ok = coord.(float64); !ok {
				return nil, ErrInvalidCoordinates
			}
		}
	case []float64:
		coords = coordinates
	case []interface{}:
		if len(coordinates) != 2 {
			return nil, ErrInvalidCoordinates
		}
		var ok bool
		if coords[0], ok = coordinates[0].(float64); !ok {
			return nil, ErrInvalidCoordinates
		}

		if coords[1], ok = coordinates[1].(float64); !ok {
			return nil, ErrInvalidCoordinates
		}
	default:
		return nil, ErrInvalidCoordinates
	}

	return coords, nil
}

// Validate validates locations.
func (driverLocation *DriverLocation) Validate() error {
	coords, err := driverLocation.getLocation()
	if err != nil {
		return err
	}

	// Check longitude limits
	if coords[0] > 180 || coords[0] < -180 {
		return fmt.Errorf("provide a valid longitude value")
	}

	// Check latitude limits
	if coords[1] > 90 || coords[1] < -90 {
		return fmt.Errorf("provide a valid latitude value")
	}

	return nil
}
