package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	{"find_nearest_param_latitude", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"latitude":191,"longitude":55,"minDistance": 0,"maxDistance": 10000}`, 400, `{"code":400,"msg":"you must provide a valid latitude"}`},
	{"find_nearest_param_longitude", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"latitude":94,"longitude":99,"minDistance": 0,"maxDistance": 10000}`, 400, `{"code":400,"msg":"you must provide a valid longitude"}`},
	{"find_nearest_param_maxDistance", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"latitude":94,"longitude":55,"minDistance": 0,"maxDistance": 0}`, 400, `{"code":400,"msg":"maxDistance must be greater then 0 and minDistance"}`},
}

var FindData = []testParams{
	{"find_nearest_success_error", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"latitude":90,"longitude":55,"minDistance": 0,"maxDistance": 10000}`, 500, `{"code":500,"msg":"Internal Server Error"}`},
}

var FindDataNotFound = []testParams{
	{"find_nearest_success_not_found", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"latitude":90,"longitude":55,"minDistance": 0,"maxDistance": 10000}`, 404, `{"code":404,"msg":"Not Found"}`},
}

var FindDataSuccess = []testParams{
	{"find_nearest_success_success", http.MethodPost, "/api/v1/driver_location/find_nearest", `{"latitude":90,"longitude":55,"minDistance": 0,"maxDistance": 10000}`, 200, `{"code":200,"data":{"locations":[{"_id":"000000000000000000000000","location":{"type":"","coordinates":null},"distance":0,"mongo_distance":0}],"total":1}}`},
}

func Test_Find_Param(t *testing.T) {

	for _, data := range FindDataParam {
		t.Run(data.name, func(t *testing.T) {
			driverLocationHandler := &httpServer{}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
			driverLocationHandler.FindNearest(w, req)

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
