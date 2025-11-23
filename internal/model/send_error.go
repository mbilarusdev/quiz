package model

import "net/http"

type SendError struct {
	W           http.ResponseWriter
	R           *http.Request
	HandlerName string
	ErrorMsg    string
	Error       error
	StatusCode  int
}
