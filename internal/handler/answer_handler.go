package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mbilarusdev/quiz/internal/model"
	"github.com/mbilarusdev/quiz/internal/service"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
	"github.com/mbilarusdev/quiz/internal/util"
)

type AnswerEndpoints interface {
	AddAnswer(
		w http.ResponseWriter,
		r *http.Request,
	)
	FindOne(
		w http.ResponseWriter,
		r *http.Request,
	)
	Delete(
		w http.ResponseWriter,
		r *http.Request,
	)
}

type AnswerHandler struct {
	AnswerSrv service.AnswerLogic
}

func NewAnswerHandler(answerSrv service.AnswerLogic) *AnswerHandler {
	res := new(AnswerHandler)
	res.AnswerSrv = answerSrv
	return res
}

func (res *AnswerHandler) AddAnswer(
	w http.ResponseWriter,
	r *http.Request,
) {
	op := "handler.AnswerHandler.AddAnswer"
	defer func() {
		if p := recover(); p != nil {
			util.SendFatal(
				model.SendFatal{
					W:           w,
					R:           r,
					HandlerName: op,
					Panic:       p,
				},
			)
		}
	}()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to parse path parameter {id}",
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to read body bytes",
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	}
	answer := &model.Answer{}
	err = json.Unmarshal(bodyBytes, answer)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed unmarshal bytes to 'Answer' model",
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	}
	if answer.QuestionID != 0 && id != answer.QuestionID {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "'Answer' field 'question_id' and path parameter 'id' not equals",
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	} else {
		answer.QuestionID = id
	}
	newAnswer, err := res.AnswerSrv.AddAnswer(answer)
	if err != nil {
		if _, ok := err.(*service_errors.NotFoundError); ok {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg: fmt.Sprintf(
						"Failed to add answer, but 'Question' with id=%v not found",
						id,
					),
					Error:      err,
					StatusCode: http.StatusUnprocessableEntity,
				},
			)
		} else if _, ok := err.(*service_errors.DuplicateError); ok {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg:    fmt.Sprintf("Failed to add answer, but 'Answer' with id=%v already exist", answer.ID),
					Error:       err,
					StatusCode:  http.StatusUnprocessableEntity,
				},
			)
		} else {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg:    "Failed to add answer, some error occured",
					Error:       err,
					StatusCode:  http.StatusUnprocessableEntity,
				},
			)
		}
		return
	}
	jsonAnswer, err := json.Marshal(newAnswer)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to marshal response with 'Answer'",
				Error:       err,
				StatusCode:  http.StatusUnprocessableEntity,
			},
		)
		return
	}

	util.SendSuccess(model.SendSuccess{
		W:           w,
		R:           r,
		HandlerName: op,
		Bytes:       jsonAnswer,
		ResultMsg: fmt.Sprintf(
			"'Answer' with id=%v for 'Question' with id=%v created succesfully",
			answer.ID,
			id,
		),
		StatusCode: http.StatusCreated,
	})
}

func (res *AnswerHandler) FindOne(
	w http.ResponseWriter,
	r *http.Request,
) {
	op := "handler.AnswerHandler.FindOne"
	defer func() {
		if p := recover(); p != nil {
			util.SendFatal(
				model.SendFatal{
					W:           w,
					R:           r,
					HandlerName: op,
					Panic:       p,
				},
			)
		}
	}()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to parse path parameter {id}",
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	}
	answer, err := res.AnswerSrv.FindOne(id)
	if err != nil {
		if _, ok := err.(*service_errors.NotFoundError); ok {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg: fmt.Sprintf(
						"'Answer' with id=%v not found!",
						id,
					),
					Error:      err,
					StatusCode: http.StatusNotFound,
				},
			)

		} else {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg: fmt.Sprintf(
						"Error when try to find 'Answer' with id=%v",
						id,
					),
					Error:      err,
					StatusCode: http.StatusUnprocessableEntity,
				},
			)
		}
		return
	}
	jsonAnswer, err := json.Marshal(answer)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to marshal response with 'Answer'",
				Error:       err,
				StatusCode:  http.StatusUnprocessableEntity,
			},
		)
		return
	}

	util.SendSuccess(model.SendSuccess{
		W:           w,
		R:           r,
		HandlerName: op,
		Bytes:       jsonAnswer,
		ResultMsg: fmt.Sprintf(
			"'Answer' with id=%v finded succesfully",
			answer.ID,
		),
		StatusCode: http.StatusOK,
	})
}

func (res *AnswerHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	op := "handler.AnswerHandler.Delete"
	defer func() {
		if p := recover(); p != nil {
			util.SendFatal(
				model.SendFatal{
					W:           w,
					R:           r,
					HandlerName: op,
					Panic:       p,
				},
			)
		}
	}()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to parse path parameter {id}",
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	}
	deleted, err := res.AnswerSrv.Delete(id)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    fmt.Sprintf("Failed to delete 'Answer' with id=%v", id),
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	}
	if !deleted {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    fmt.Sprintf("'Answer' with id=%v not found!", id),
				Error:       err,
				StatusCode:  http.StatusNotFound,
			},
		)
		return
	}

	util.SendSuccess(model.SendSuccess{
		W:           w,
		R:           r,
		HandlerName: op,
		Bytes:       make([]byte, 0),
		ResultMsg: fmt.Sprintf(
			"'Answer' with id=%v deleted succesfully",
			id,
		),
		StatusCode: http.StatusNoContent,
	})
}
