package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/s3f4/locationmatcher/internal/driverlocation/mocks"
	"github.com/s3f4/locationmatcher/internal/driverlocation/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testParams struct {
	name         string
	method       string
	url          string
	body         string
	expectedCode int
	expectedBody string
}

var FindDataParam = []testParams{
	{"find_nearest_param_no_param", http.MethodPost, "/api/v1/driver_location/find_nearest", ``, 400, `{"code":400,"msg":"Bad Request"}`},
	{"find_nearest_param_invalid_geojson", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"latitude":191,"longitude":55,"minDistance": 0,"maxDistance": 10000}`, 400, `{"code":400,"msg":"you must provide a valid GeoJSON type"}`},
	{"find_nearest_param_invalid_type", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"location": {"type": "Poin","coordinates": [29.15188821,41.90513187]},"minDistance": 55,"maxDistance": 10000}`, 400, `{"code":400,"msg":"you must provide a valid GeoJSON type"}`},
	{"find_nearest_param_invalid_longitude", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"location": {"type": "Point","coordinates": [181,29.15188821]},"minDistance": 55,"maxDistance": 10000}`, 400, `{"code":400,"msg":"you must provide a valid longitude"}`},
	{"find_nearest_param_invalid_latitude", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"location": {"type": "Point","coordinates": [41.90513187,-91]},"minDistance": 55,"maxDistance": 10000}`, 400, `{"code":400,"msg":"you must provide a valid latitude"}`},
	{"find_nearest_param_maxDistance", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"location": {"type": "Point","coordinates": [41.90513187,29.15188821]},"minDistance": 55,"maxDistance": 0}`, 400, `{"code":400,"msg":"maxDistance must be greater then 0 and minDistance"}`},
}

var FindData = []testParams{
	{"find_nearest_success_error", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"location": {"type": "Point","coordinates": [41.90513187,29.15188821]},"minDistance": 55,"maxDistance": 10000}`, 500, `{"code":500,"msg":"Internal Server Error"}`},
}

var FindDataNotFound = []testParams{
	{"find_nearest_success_not_found", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"location": {"type": "Point","coordinates": [41.90513187,29.15188821]},"minDistance": 55,"maxDistance": 10000}`, 404, `{"code":404,"msg":"Not Found"}`},
}

var FindDataSuccess = []testParams{
	{"find_nearest_success_success", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"location": {"type": "Point","coordinates": [41.90513187,29.15188821]},"minDistance": 55,"maxDistance": 10000}`, 200, `{"code":200,"data":{"total":1,"locations":[{"_id":"000000000000000000000000","location":{"type":"","coordinates":null},"distance":0}]}}`},
}

var UpsertBulkParams = []testParams{
	{"upsertbulk_", http.MethodPost, "/api/v1/driver_location", ``, 400, `{"code":400,"msg":"Bad Request"}`},
	{"upsertbulk_parse_error", http.MethodPost, "/api/v1/driver_location", `[]`, 400, `{"code":400,"msg":"provide valid driver locations"}`},
	{"upsertbulk_invalid_latitude", http.MethodPost, "/api/v1/driver_location", `[{"_id":"6219f72c61d60d9a30ff2072","location":{"type":"Point","coordinates":[-190.94001079,29.00077262]}}]`, 400, `{"code":400,"msg":"provide valid driver locations"}`},
	{"upsertbulk_invalid_longitude", http.MethodPost, "/api/v1/driver_location", `[{"_id":"6219f72c61d60d9a30ff2072","location":{"type":"Point","coordinates":[40.94001079,191.00077262]}}]`, 400, `{"code":400,"msg":"provide valid driver locations"}`},
	// {"driver_location_valid_request", http.MethodPost, "/api/v1/driver_location", `[{"_id":"6219f72c61d60d9a30ff2072","location":{"type":"Point","coordinates":[40.94001079,29.00077262]}}]`, 200, `[{"_id":"6219f72c61d60d9a30ff2072","location":{"type":"Point","coordinates":[40.94001079,29.00077262]}}]`},
}

var UpsertBulkValues = []testParams{
	{"driver_location_valid_request", http.MethodPost, "/api/v1/driver_location", `[{"_id":"6219f72c61d60d9a30ff2072","location":{"type":"Point","coordinates":[40.94001079,29.00077262]}}]`, 200, `[{"_id":"6219f72c61d60d9a30ff2072","location":{"type":"Point","coordinates":[40.94001079,29.00077262]},"distance":0}]`},
}

var UpsertBulkErr = []testParams{
	{"driver_location_valid_request_error", http.MethodPost, "/api/v1/driver_location", `[{"_id":"6219f72c61d60d9a30ff2072","location":{"type":"Point","coordinates":[40.94001079,29.00077262]}}]`, 500, `{"code":500,"msg":"Internal Server Error"}`},
}

func Test_Find_Param(t *testing.T) {
	driverLocationRepository := new(mocks.Repository)
	driverLocationRepository.On("Find", context.TODO(), mock.Anything).Return(nil, errors.New("error"))

	for _, data := range FindDataParam {
		t.Run(data.name, func(t *testing.T) {
			driverLocationHandler := &httpServer{driverLocationRepository}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
			driverLocationHandler.Find(w, req)

			res := w.Result()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, data.expectedBody, string(body))
			assert.Equal(t, data.expectedCode, w.Code)
		})
	}
}

func Test_Find(t *testing.T) {
	driverLocationRepository := new(mocks.Repository)
	driverLocationRepository.On("Find1", context.TODO(), &models.Query{
		Location: models.Location{
			Type:        "Point",
			Coordinates: []interface{}{41.90513187, 29.15188821},
		},
		MinDistance: 55,
		MaxDistance: 10000,
	}).Return(nil, errors.New("error"))

	for _, data := range FindData {
		t.Run(data.name, func(t *testing.T) {
			driverLocationHandler := &httpServer{driverLocationRepository}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
			driverLocationHandler.Find(w, req)

			res := w.Result()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, data.expectedBody, string(body))
			assert.Equal(t, data.expectedCode, w.Code)
		})

	}
}

func Test_Find_NotFound(t *testing.T) {
	driverLocationRepository := new(mocks.Repository)
	driverLocationRepository.On("Find1", context.TODO(), &models.Query{
		Location: models.Location{
			Type:        "Point",
			Coordinates: []interface{}{41.90513187, 29.15188821},
		},
		MinDistance: 55,
		MaxDistance: 10000,
	}).Return([]*models.DriverLocation{}, nil)

	for _, data := range FindDataNotFound {
		t.Run(data.name, func(t *testing.T) {
			driverLocationHandler := &httpServer{driverLocationRepository}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
			driverLocationHandler.Find(w, req)

			res := w.Result()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, data.expectedBody, string(body))
			assert.Equal(t, data.expectedCode, w.Code)
		})
	}
}

func Test_Find_Success(t *testing.T) {
	driverLocationRepository := new(mocks.Repository)
	driverLocationRepository.On("Find1", context.TODO(), &models.Query{
		Location: models.Location{
			Type:        "Point",
			Coordinates: []interface{}{41.90513187, 29.15188821},
		},
		MinDistance: 55,
		MaxDistance: 10000,
	}).Return([]*models.DriverLocation{{Location: models.Location{}}}, nil)

	for _, data := range FindDataSuccess {
		t.Run(data.name, func(t *testing.T) {
			driverLocationHandler := &httpServer{driverLocationRepository}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
			driverLocationHandler.Find(w, req)

			res := w.Result()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}

			fmt.Println(string(body))
			assert.Equal(t, data.expectedBody, string(body))
			assert.Equal(t, data.expectedCode, w.Code)
		})
	}
}

func Test_UpsertBulk_Params(t *testing.T) {
	driverLocationRepository := new(mocks.Repository)
	driverLocationRepository.On("UpsertBulk", context.TODO(), []*models.DriverLocation{}).Return(nil)

	for _, data := range UpsertBulkParams {
		driverLocationHandler := &httpServer{driverLocationRepository}
		w := httptest.NewRecorder()
		req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
		driverLocationHandler.UpsertBulk(w, req)

		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, data.expectedBody, string(body))
		assert.Equal(t, data.expectedCode, w.Code)
	}
}

func Test_UpsertBulk_Valid(t *testing.T) {
	driverLocationRepository := new(mocks.Repository)
	id, _ := primitive.ObjectIDFromHex("6219f72c61d60d9a30ff2072")
	driverLocationRepository.On("UpsertBulk", context.TODO(), []*models.DriverLocation{
		{
			ID: id,
			Location: models.Location{
				Type:        "Point",
				Coordinates: []interface{}{40.94001079, 29.00077262},
			},
		},
	}).Return(nil)

	for _, data := range UpsertBulkValues {
		driverLocationHandler := &httpServer{driverLocationRepository}
		w := httptest.NewRecorder()
		req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
		driverLocationHandler.UpsertBulk(w, req)

		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, data.expectedBody, string(body))
		assert.Equal(t, data.expectedCode, w.Code)
	}
}

func Test_UpsertBulk_Error(t *testing.T) {
	driverLocationRepository := new(mocks.Repository)
	id, _ := primitive.ObjectIDFromHex("6219f72c61d60d9a30ff2072")
	driverLocationRepository.On("UpsertBulk", context.TODO(), []*models.DriverLocation{
		{
			ID: id,
			Location: models.Location{
				Type:        "Point",
				Coordinates: []interface{}{40.94001079, 29.00077262},
			},
		},
	}).Return(fmt.Errorf("err"))

	for _, data := range UpsertBulkErr {
		driverLocationHandler := &httpServer{driverLocationRepository}
		w := httptest.NewRecorder()
		req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
		driverLocationHandler.UpsertBulk(w, req)

		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, data.expectedBody, string(body))
		assert.Equal(t, data.expectedCode, w.Code)
	}
}
