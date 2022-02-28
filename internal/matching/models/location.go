package models

type Location struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}
