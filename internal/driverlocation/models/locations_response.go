package models

type LocationsResponse struct {
	Total     int               `json:"total"`
	Locations []*DriverLocation `json:"locations"`
}
