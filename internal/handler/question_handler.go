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
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		util.SendError(w, "Failed read body bytes\n", http.StatusBadRequest)

		return
	}
	question := &model.Question{}
	err = json.Unmarshal(bodyBytes, question)
	if err != nil {
		util.SendError(w, "Failed unmarshal bytes to question\n", http.StatusBadRequest)
		return
	}
	newQuestion, err := res.QuestionSrv.Create(question)
	if err != nil {
		if _, ok := err.(*service_errors.DuplicateError); ok {
			util.SendError(
				w,
				fmt.Sprintf(
					"Failed to create question with id=%v, but already exist\n",
					question.ID,
				),
				http.StatusUnprocessableEntity,
			)
		} else {
			util.SendError(w, "Failed to create question\n", http.StatusUnprocessableEntity)
		}
		return
	}
	jsonQuestion, err := json.Marshal(newQuestion)
	if err != nil {
		util.SendError(w, "Failed to marshal question response\n", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Question with id=%vcreated succesfully\n", newQuestion.ID)
	util.SendSuccess(w, jsonQuestion, http.StatusCreated)
}

func (res *QuestionHandler) FindOneDetailed(
	w http.ResponseWriter,
	r *http.Request,
) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(w, "Failed to parse question field id\n", http.StatusBadRequest)
		return
	}
	question, err := res.QuestionSrv.FindOneDetailed(id)
	if err != nil {
		if _, ok := err.(*service_errors.NotFoundError); ok {
			util.SendError(
				w,
				fmt.Sprintf("Question with id=%v not found\n", id),
				http.StatusNotFound,
			)
		} else {
			util.SendError(
				w,
				fmt.Sprintf("Failed to find question with id=%v and his answers\n", id),
				http.StatusUnprocessableEntity,
			)
		}
		return
	}
	jsonQuestion, err := json.Marshal(question)
	if err != nil {
		util.SendError(w, "Failed to marshal question response\n", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Question with id=%v and his answers finded successfully\n", id)
	util.SendSuccess(w, jsonQuestion, http.StatusOK)
}

func (res *QuestionHandler) FindAll(
	w http.ResponseWriter,
	r *http.Request,
) {

	questions, err := res.QuestionSrv.FindAll()
	if err != nil {
		util.SendError(
			w,
			"Failed to find all answers\n",
			http.StatusNotFound,
		)
		return
	}
	jsonQuestions, err := json.Marshal(questions)
	if err != nil {
		util.SendError(
			w,
			"Failed to marshal find all questions response\n",
			http.StatusInternalServerError,
		)
		return
	}

	fmt.Printf("Find all questions successfully\n")
	util.SendSuccess(w, jsonQuestions, http.StatusOK)
}

func (res *QuestionHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(w, "Failed to parse question field id\n", http.StatusBadRequest)
		return
	}
	_, err = res.QuestionSrv.Delete(id)
	if err != nil {
		util.SendError(
			w,
			fmt.Sprintf("Failed to delete question with id=%v\n", id),
			http.StatusNotFound,
		)
		return
	}

	fmt.Printf("Question with id=%v deleted successfully\n", id)
	util.SendSuccess(w, make([]byte, 0), http.StatusNoContent)
}
