//  The Matching Api
//   version: 0.0.1
//   title: Driver Location Api
//  Schemes: http, https
//  Host: localhost:3001
//  BasePath: /api/v1
//  Produces:
//    - application/json
//
// securityDefinitions:
//  Bearer:
//    type: apiKey
//    in: header
//    name: Authorization
// swagger:meta
package server

type Location struct {
	// example: Point
	Type string `json:"type"`
	// example: [41.90513187,29.15188821]
	Coordinates [2]float64 `json:"coordinates" example:"Point"`
}

// swagger:model
type LocationBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/Location"
	//  required: true
	Body Location `json:"body"`
}

// swagger:model DriverLocation
type DriverLocation struct {
	// Id of the driver location
	// in: string
	ID string `json:"_id"`
	// Id of the driver location
	// in: Location
	Location Location `json:"location"`
	// Distance of the driver location to given coordinates
	// readonly: true
	Distance float64 `json:"distance"`
	// Distance of the driver location to given coordinates that comes from mongodb
	// readonly: true
	MongoDistance *float64 `json:"mongo_distance"`
}

type Query struct {
	// Id of the driver location
	// in: Location
	Location Location `json:"location"`
	// Minimum distance in meters
	// in: float64
	MinDistance int64 `json:"minDistance"`
	// Maximum distance in meters
	// in: float64
	// example: 10000
	MaxDistance int64 `json:"maxDistance"`
}

// swagger:parameters v1 Find
type QueryBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/Query"
	//  required: true
	Body Query `json:"body"`
}

type DriverLocations struct {
	DriverLocations []*DriverLocation
}

// swagger:parameters v1 UpsertBulk
type ReqUpsertBulk struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: array
	//     "$ref": "#/definitions/DriverLocations"
	//  required: true
	Body []*DriverLocation `json:"body"`
}

type ApiError struct {
	// Status Code of the error
	// in: int
	Code int
	// Message of the error
	// in: string
	Msg string
}

// swagger:response ApiError
type ApiErrBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/responses/ApiError"
	//  required: true
	Body ApiError `json:"body"`
}

type Response struct {
	// Status Code of the error
	// in: int
	Code int `json:"code"`

	// Message of the error
	// in: string
	Msg string `json:"msg"`
	// Message of the error
	// in: string
	Data string `json:"data"`
}

// swagger:response Response
type ResponseP struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/responses/Response"
	//  required: true
	Body Response `json:"body"`
}
