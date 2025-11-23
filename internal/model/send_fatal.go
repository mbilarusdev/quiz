package model

import "net/http"

type SendFatal struct {
	W           http.ResponseWriter
	R           *http.Request
	HandlerName string
	Panic       any
}
