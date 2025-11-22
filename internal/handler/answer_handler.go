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
	FindAllByQuestionID(
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
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(w, "Failed to parse question field id\n", http.StatusBadRequest)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		util.SendError(w, "Failed read body bytes\n", http.StatusBadRequest)

		return
	}
	answer := &model.Answer{}
	err = json.Unmarshal(bodyBytes, answer)
	if err != nil {
		util.SendError(w, "Failed unmarshal bytes to answer\n", http.StatusBadRequest)
		return
	}
	if answer.QuestionID != 0 && id != answer.QuestionID {
		util.SendError(
			w,
			fmt.Sprintf(
				"Failed add answer with id=%v to question with id=%v, but question id in path parameters and body not equals\n",
				answer.ID,
				id,
			),
			http.StatusBadRequest,
		)
		return
	}
	answer.QuestionID = id
	newAnswer, err := res.AnswerSrv.AddAnswer(answer)
	if err != nil {
		if _, ok := err.(*service_errors.NotFoundError); ok {
			util.SendError(
				w,
				fmt.Sprintf(
					"Failed add answer with id=%v to question with id=%v, but question not found\n",
					answer.ID,
					id,
				),
				http.StatusNotFound,
			)
		} else if _, ok := err.(*service_errors.DuplicateError); ok {
			util.SendError(
				w,
				fmt.Sprintf("Failed add answer with id=%v to question with id=%v, but answer with id=%v already exist\n", answer.ID, id, answer.ID),
				http.StatusUnprocessableEntity,
			)
		} else {
			util.SendError(
				w,
				fmt.Sprintf("Failed add answer with id=%v to question with id=%v\n", answer.ID, id),
				http.StatusUnprocessableEntity,
			)
		}
		return
	}
	jsonAnswer, err := json.Marshal(newAnswer)
	if err != nil {
		util.SendError(w, "Failed to marshal answer response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Answer with id=%v for question with id=%v created succesfully\n", answer.ID, id)
	util.SendSuccess(w, jsonAnswer, http.StatusCreated)
}

func (res *AnswerHandler) FindAllByQuestionID(
	w http.ResponseWriter,
	r *http.Request,
) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(w, "Failed to parse question field id", http.StatusBadRequest)
		return
	}
	answers, err := res.AnswerSrv.FindAllByQuestionID(id)
	if err != nil {
		util.SendError(
			w,
			fmt.Sprintf("Failed to find answers by question with id=%v\n", id),
			http.StatusNotFound,
		)
		return
	}
	jsonAnswers, err := json.Marshal(answers)
	if err != nil {
		util.SendError(w, "Failed to marshal answers response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Answers for question with id=%v finded successfully\n", id)
	util.SendSuccess(w, jsonAnswers, http.StatusOK)
}

func (res *AnswerHandler) FindOne(
	w http.ResponseWriter,
	r *http.Request,
) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(w, "Failed to parse answer field id", http.StatusBadRequest)
		return
	}
	question, err := res.AnswerSrv.FindOne(id)
	if err != nil {
		if _, ok := err.(*service_errors.NotFoundError); ok {
			util.SendError(
				w,
				fmt.Sprintf("Answer with id=%v not found\n", id),
				http.StatusNotFound,
			)

		} else {
			util.SendError(
				w,
				fmt.Sprintf("Failed to find answer with id=%v\n", id),
				http.StatusUnprocessableEntity,
			)
		}
		return
	}
	jsonQuestion, err := json.Marshal(question)
	if err != nil {
		util.SendError(w, "Failed to marshal question response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Question with id=%v and his answers finded successfully\n", id)
	util.SendSuccess(w, jsonQuestion, http.StatusOK)
}

func (res *AnswerHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.SendError(w, "Failed to parse answer field id", http.StatusBadRequest)
		return
	}
	_, err = res.AnswerSrv.Delete(id)
	if err != nil {
		util.SendError(
			w,
			fmt.Sprintf("Failed to delete answer with id=%v\n", id),
			http.StatusNotFound,
		)
		return
	}

	fmt.Printf("Answer with id=%v deleted successfully\n", id)
	util.SendSuccess(w, make([]byte, 0), http.StatusNoContent)
}
