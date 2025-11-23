package model

import "net/http"

type SendSuccess struct {
	W           http.ResponseWriter
	R           *http.Request
	HandlerName string
	Bytes       []byte
	ResultMsg   string
	StatusCode  int
}
