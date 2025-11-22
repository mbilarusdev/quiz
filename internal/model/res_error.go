package model

type ResError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}
