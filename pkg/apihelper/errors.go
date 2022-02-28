package apihelper

import "net/http"

// ApiError holds errors to show to the user.
type ApiError struct {
	Code int
	Msg  string
}

var _ error = (*ApiError)(nil)

func NewApiError(code int, msg string) *ApiError {
	return &ApiError{
		Code: code,
		Msg:  msg,
	}
}

func (e *ApiError) Error() string {
	return e.Msg
}

var (
	Err400 = NewApiError(http.StatusBadRequest, "Bad Request")
	Err422 = NewApiError(http.StatusUnprocessableEntity, "Unprocessable Entity")
	Err401 = NewApiError(http.StatusUnauthorized, "Unauthorized")
	Err403 = NewApiError(http.StatusForbidden, "Forbidden")
	Err404 = NewApiError(http.StatusNotFound, "Not Found")
	Err500 = NewApiError(http.StatusInternalServerError, "Internal Server Error")
)
