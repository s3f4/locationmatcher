package models

import (
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_CalculateDistance(t *testing.T) {
	// calculate distance between the Statue of liberty and the Eiffel tower
	// 5837.41km
	// Statue of liberty
	driverLocation := DriverLocation{
		Location: Location{
			Type:        "Point",
			Coordinates: []float64{40.6892, -74.0444},
		},
	}

	// Eiffel
	d, err := driverLocation.CalculateDistance(48.8583, 2.2945)
	if err != nil {
		t.Error("Wrong coordinates")
	}

	if fmt.Sprintf("%.2f", d) != fmt.Sprintf("%.2f", 5837.41) {
		fmt.Println(d)
		t.Error("Long distance calculation error")
	}

	// calculate distance between the Galata tower and Ayasofya
	// 1.96km
	// galata
	driverLocation = DriverLocation{
		Location: Location{
			Type:        "Point",
			Coordinates: []float64{41.025651081666744, 28.97413088610361},
		},
	}

	// ayasofya
	d, err = driverLocation.CalculateDistance(41.00858654897259, 28.979986854317975)
	if err != nil {
		t.Error("Wrong coordinates")
	}
	if fmt.Sprintf("%.2f", d) != fmt.Sprintf("%.2f", 1.96) {
		fmt.Println(d)
		t.Error("Short distance calculation error")
	}

	// error conditions
	driverLocation = DriverLocation{
		Location: Location{
			Type:        "Point",
			Coordinates: []int{41, 28},
		},
	}

	_, err = driverLocation.CalculateDistance(41.00858654897259, 28.979986854317975)
	if err == nil {
		t.Error("It must be return an error")
	}
}

func Test_getLocation(t *testing.T) {
	driverLocation := DriverLocation{
		Location: Location{
			Coordinates: []float64{40.94, 29.1},
		},
	}

	coords, err := driverLocation.getLocation()
	if err != nil {
		t.Error(err)
	}

	if coords[0] != 40.94 || coords[1] != 29.1 {
		t.Error("wrong coordinates")
	}

	if reflect.TypeOf(coords) != reflect.TypeOf([]float64{}) {
		t.Error("wrong coordinate types")
	}

	driverLocation.Location.Coordinates = primitive.A{40.94, 29.1}
	coords, err = driverLocation.getLocation()
	if err != nil {
		t.Error(err)
	}

	if coords[0] != 40.94 || coords[1] != 29.1 {
		t.Error("wrong coordinates")
	}

	if reflect.TypeOf(coords) != reflect.TypeOf([]float64{}) {
		t.Error("wrong coordinate types")
	}

	driverLocation.Location.Coordinates = []interface{}{40}
	if _, err = driverLocation.getLocation(); err == nil {
		t.Error("wrong coordinate types")
	}

	driverLocation.Location.Coordinates = []interface{}{40.22, int64(51)}
	if _, err = driverLocation.getLocation(); err == nil {
		t.Error("wrong coordinate types")
	}

	driverLocation.Location.Coordinates = []interface{}{40, 29}
	if _, err = driverLocation.getLocation(); err == nil {
		t.Error("wrong coordinate types")
	}

	driverLocation.Location.Coordinates = []int{40, 29}
	if _, err = driverLocation.getLocation(); err == nil {
		t.Error("wrong coordinate types")
	}

	driverLocation.Location.Coordinates = primitive.A{40, 29}
	if _, err = driverLocation.getLocation(); err == nil {
		t.Error("wrong coordinate types")
	}
}

func Test_Validate(t *testing.T) {
	driverLocation := DriverLocation{
		Location: Location{
			Coordinates: []float64{40.94, 29.1},
		},
	}

	if err := driverLocation.Validate(); err != nil {
		t.Error("It must not return an error")
	}

	driverLocation.Location.Coordinates.([]float64)[0] = -181.22
	if err := driverLocation.Validate(); err != nil {
		if err.Error() != "provide a valid latitude value" {
			t.Error("Wrong latitude validation message")
		}
	}

	driverLocation.Location.Coordinates.([]float64)[0] = 0
	driverLocation.Location.Coordinates.([]float64)[1] = 91
	if err := driverLocation.Validate(); err != nil {
		if err.Error() != "provide a valid longitude value" {
			t.Error("Wrong longitude validation message")
		}
	}

	driverLocation.Location.Coordinates = []int{1, 2}
	if err := driverLocation.Validate(); err != nil {
		if err.Error() != "invalid coordinates" {
			t.Error("Coordinate error")
		}
	}
}
