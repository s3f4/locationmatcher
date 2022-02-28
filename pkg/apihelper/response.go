package apihelper

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response holds response data
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

//SendResponse returns json response
func SendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp, err := json.Marshal(data)
	if err != nil {
		log.Println("Body JSON Marshal error")
		return
	}

	w.Write(resp)
}

func Send400(w http.ResponseWriter) {
	SendResponse(
		w,
		Err400.Code, Response{
			Code: Err400.Code,
			Msg:  Err400.Error(),
		})
}

func Send401(w http.ResponseWriter) {
	SendResponse(
		w,
		Err400.Code, Response{
			Code: Err400.Code,
			Msg:  Err400.Error(),
		})
}

// Send404 sends not found error
func Send404(w http.ResponseWriter) {
	SendResponse(
		w,
		Err404.Code, Response{
			Code: Err404.Code,
			Msg:  Err404.Error(),
		})
}

// Send500 sends internal server error
func Send500(w http.ResponseWriter) {
	SendResponse(
		w,
		Err500.Code,
		Response{
			Code: Err500.Code,
			Msg:  Err500.Error(),
		},
	)
}
