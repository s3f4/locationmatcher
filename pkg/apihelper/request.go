package apihelper

import (
	"encoding/json"
	"net/http"
)

type RequestBody interface {
	Validate() error
}

// ParseAndValidate parses json body and makes validation on struct
func ParseAndValidate(r *http.Request, model RequestBody) *ApiError {
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		return Err400
	}
	if err := model.Validate(); err != nil {
		return NewApiError(400, err.Error())
	}

	return nil
}
