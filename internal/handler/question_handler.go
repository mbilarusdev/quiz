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

type QuestionEndpoints interface {
	AddAnswer(
		w http.ResponseWriter,
		r *http.Request,
	)
	FindOneDetailed(
		w http.ResponseWriter,
		r *http.Request,
	)
	FindAll(
		w http.ResponseWriter,
		r *http.Request,
	)
	Delete(
		w http.ResponseWriter,
		r *http.Request,
	)
}

type QuestionHandler struct {
	QuestionSrv service.QuestionLogic
}

func NewQuestionHandler(questionSrv service.QuestionLogic) *QuestionHandler {
	res := new(QuestionHandler)
	res.QuestionSrv = questionSrv
	return res
}

func (res *QuestionHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	op := "handler.QuestionHandler.Create"
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
	question := &model.Question{}
	err = json.Unmarshal(bodyBytes, question)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed unmarshal bytes to 'Question' model",
				Error:       err,
				StatusCode:  http.StatusBadRequest,
			},
		)
		return
	}
	newQuestion, err := res.QuestionSrv.Create(question)
	if err != nil {
		if _, ok := err.(*service_errors.DuplicateError); ok {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg: fmt.Sprintf(
						"Failed to create question, but 'Question' with id=%v already exist",
						question.ID,
					),
					Error:      err,
					StatusCode: http.StatusUnprocessableEntity,
				},
			)
		} else {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg:    "Failed to create question, some error occured",

					Error:      err,
					StatusCode: http.StatusUnprocessableEntity,
				},
			)
		}
		return
	}
	jsonQuestion, err := json.Marshal(newQuestion)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to marshal response with 'Question'",
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
		Bytes:       jsonQuestion,
		ResultMsg: fmt.Sprintf(
			"'Question' with id=%v created succesfully",
			newQuestion.ID,
		),
		StatusCode: http.StatusCreated,
	})
}

func (res *QuestionHandler) FindOneDetailed(
	w http.ResponseWriter,
	r *http.Request,
) {
	op := "handler.QuestionHandler.FindOneDetailed"
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
	question, err := res.QuestionSrv.FindOneDetailed(id)
	if err != nil {
		if _, ok := err.(*service_errors.NotFoundError); ok {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg: fmt.Sprintf(
						"'Question' with id=%v not found",
						id,
					),
					Error:      err,
					StatusCode: http.StatusUnprocessableEntity,
				},
			)
		} else {
			util.SendError(
				model.SendError{
					W:           w,
					R:           r,
					HandlerName: op,
					ErrorMsg: fmt.Sprintf(
						"Failed to find 'Question' with id=%v, some error occured",
						id,
					),
					Error:      err,
					StatusCode: http.StatusUnprocessableEntity,
				},
			)
		}
		return
	}
	jsonQuestion, err := json.Marshal(question)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to marshal response with 'Question'",
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
		Bytes:       jsonQuestion,
		ResultMsg: fmt.Sprintf(
			"'Question' with id=%v and his answers finded succesfully",
			id,
		),
		StatusCode: http.StatusOK,
	})
}

func (res *QuestionHandler) FindAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	op := "handler.QuestionHandler.FindAll"
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
	questions, err := res.QuestionSrv.FindAll()
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to find all questions, some error occured",
				Error:       err,
				StatusCode:  http.StatusUnprocessableEntity,
			},
		)
		return
	}
	jsonQuestions, err := json.Marshal(questions)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg:    "Failed to marshal response with list of 'Question'",
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
		Bytes:       jsonQuestions,
		ResultMsg:   "List of 'Question' finded succesfully",
		StatusCode:  http.StatusOK,
	})
}

func (res *QuestionHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	op := "handler.QuestionHandler.Delete"
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
	deleted, err := res.QuestionSrv.Delete(id)
	if err != nil {
		util.SendError(
			model.SendError{
				W:           w,
				R:           r,
				HandlerName: op,
				ErrorMsg: fmt.Sprintf(
					"Failed to delete 'Question' with id=%v, some error occured",
					id,
				),
				Error:      err,
				StatusCode: http.StatusUnprocessableEntity,
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
				ErrorMsg:    fmt.Sprintf("'Question' with id=%v not found!", id),
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
			"'Question' with id=%v deleted succesfully",
			id,
		),
		StatusCode: http.StatusNoContent,
	})
}
