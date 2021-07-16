package controller

import (
	"encoding/json"
	"log"

	. "github.com/zackwn/article-api/usecase"
)

func StatusResponse(statusCode int) *Response {
	return &Response{StatusCode: statusCode}
}

func ErrorResponse(useCaseErr UseCaseErr) *Response {
	errRes := struct {
		Message string `json:"error"`
	}{Message: useCaseErr.Error()}
	json, err := json.Marshal(&errRes)
	if err != nil {
		log.Fatal(useCaseErr)
	}
	return &Response{StatusCode: useCaseErr.HttpStatus(), Data: json}
}

// JSONResponse takes the status code and the pointer to a response struct
func JSONResponse(statusCode int, structRes interface{}) *Response {
	json, err := json.Marshal(structRes)
	if err != nil {
		log.Fatal(err)
	}
	return &Response{StatusCode: statusCode, Data: json}
}
