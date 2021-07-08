package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

func StatusResponse(statusCode int) *Response {
	return &Response{StatusCode: statusCode}
}

func ErrorResponse(err error) *Response {
	errRes := struct {
		Message string `json:"error"`
	}{Message: err.Error()}
	json, err := json.Marshal(&errRes)
	if err != nil {
		log.Fatal(err)
	}
	return &Response{StatusCode: http.StatusInternalServerError, Data: json}
}

// JSONResponse takes the status code and the pointer to a response struct
func JSONResponse(statusCode int, structRes interface{}) *Response {
	json, err := json.Marshal(structRes)
	if err != nil {
		log.Fatal(err)
	}
	return &Response{StatusCode: statusCode, Data: json}
}
