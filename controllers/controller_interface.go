package controller

import "net/http"

type Response struct {
	StatusCode int
	Data       []byte
}

type Controller interface {
	Handle(request *http.Request) *Response
}
