package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbilarusdev/quiz/internal/model"
)

func SendError(w http.ResponseWriter, errorMsg string, statusCode int) {
	responseError := model.ResError{Error: errorMsg, StatusCode: statusCode}
	fmt.Printf("%s", responseError.Error)
	jsonErr, err := json.Marshal(responseError)
	if err != nil {
		fmt.Printf("Error when marshal error response\n")
	}
	w.WriteHeader(statusCode)
	w.Write(jsonErr)
}

func SendSuccess(w http.ResponseWriter, data []byte, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write(data)
}
