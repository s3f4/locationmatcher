package models

import (
	"testing"
)

func Test_Query(t *testing.T) {
	tests := []struct {
		name    string
		query   Query
		wantErr bool
	}{
		{
			name: "it should return an error if type is not valid",
			query: Query{
				Location: Location{
					Type:        "Poin",
					Coordinates: []interface{}{10.0, 10.0},
				},
				MinDistance: 0,
				MaxDistance: 121,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if coords is not valid",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []float64{10.0, 10.0},
				},
				MinDistance: 0,
				MaxDistance: 121,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if coords len greater than 2",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{10.0, 10.0, 10.0, 10.0},
				},
				MinDistance: 0,
				MaxDistance: 121,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if latitude is not float",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{10, 10.0},
				},
				MinDistance: 0,
				MaxDistance: 121,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if latitude is not valid",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{191.1, 10.0},
				},
				MinDistance: 0,
				MaxDistance: 121,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if longitude is not float",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{10.0, 10},
				},
				MinDistance: 0,
				MaxDistance: 10,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if longitude is not valid",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{10.0, 91.1},
				},
				MinDistance: 0,
				MaxDistance: 10,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if minDistance >= maxDistance",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{10, 10},
				},
				MinDistance: 102,
				MaxDistance: 102,
			},
			wantErr: true,
		},
		{
			name: "it should return an error if maxDistance = 0",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{10.0, 10.0},
				},
				MinDistance: 0,
				MaxDistance: 0,
			},
			wantErr: true,
		},
		{
			name: "it shouldn't return an error ",
			query: Query{
				Location: Location{
					Type:        "Point",
					Coordinates: []interface{}{10.1, 10.0},
				},
				MinDistance: 0,
				MaxDistance: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.query.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
