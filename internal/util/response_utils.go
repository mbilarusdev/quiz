package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbilarusdev/quiz/internal/common"
	"github.com/mbilarusdev/quiz/internal/model"
	"go.uber.org/zap"
)

func SendError(
	params model.SendError,
) {
	responseError := model.ResError{Error: params.ErrorMsg, StatusCode: params.StatusCode}
	common.L.Error(
		fmt.Sprintf("Handler '%v' returns error", params.HandlerName),
		zap.String("Path", params.R.URL.Path),
		zap.String("Method", params.R.Method),
		zap.String("ErrorMsg", params.ErrorMsg),
		zap.String("Error", params.Error.Error()),
		zap.Int("StatusCode", params.StatusCode),
	)
	jsonErr, err := json.Marshal(responseError)
	if err != nil {
		common.L.Fatal("json.Marshal failed in util.SendError")
	}
	params.W.WriteHeader(params.StatusCode)
	params.W.Write(jsonErr)
}

func SendSuccess(
	params model.SendSuccess,
) {
	common.L.Info(
		fmt.Sprintf("Handler '%v' returns success", params.HandlerName),
		zap.String("Path", params.R.URL.Path),
		zap.String("Method", params.R.Method),
		zap.String("Result", params.ResultMsg),
		zap.Int("StatusCode", params.StatusCode),
	)
	params.W.WriteHeader(params.StatusCode)
	params.W.Write(params.Bytes)
}

func SendFatal(
	params model.SendFatal,
) {
	statusCode := http.StatusInternalServerError
	responseError := model.ResError{Error: "Panic occured on server side", StatusCode: statusCode}
	common.L.Error(
		fmt.Sprintf("Handler '%v' occured with panic", params.HandlerName),
		zap.String("Path", params.R.URL.Path),
		zap.String("Method", params.R.Method),
		zap.Any("Panic", params.Panic),
		zap.Int("StatusCode", statusCode),
	)
	jsonErr, err := json.Marshal(responseError)
	if err != nil {
		common.L.Error("json.Marshal failed in util.SendFatal")
	}
	params.W.WriteHeader(statusCode)
	params.W.Write(jsonErr)
}
