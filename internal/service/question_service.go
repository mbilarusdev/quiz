package service

import (
	"context"

	"github.com/mbilarusdev/quiz/internal/model"
	"github.com/mbilarusdev/quiz/internal/repository"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
)

type QuestionLogic interface {
	Create(question *model.Question) (*model.Question, error)
	FindOneDetailed(questionID int) (*model.Question, error)
	FindAll() ([]model.Question, error)
	Delete(questionID int) (bool, error)
}

type QuestionService struct {
	QuestionRepo repository.QuestionProvider
	AnswerRepo   repository.AnswerProvider
}

func NewQuestionService(
	questionRepo repository.QuestionProvider,
	answerRepo repository.AnswerProvider,
) *QuestionService {
	service := new(QuestionService)
	service.QuestionRepo = questionRepo
	service.AnswerRepo = answerRepo
	return service
}

func (service *QuestionService) Create(question *model.Question) (*model.Question, error) {
	question, err := service.QuestionRepo.Insert(context.Background(), question)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (service *QuestionService) FindOneDetailed(questionID int) (*model.Question, error) {
	ctx := context.Background()
	question, err := service.QuestionRepo.GetOne(nil, ctx, questionID, true)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return question, &service_errors.NotFoundError{ID: questionID}
	}

	return question, nil
}

func (service *QuestionService) FindAll() ([]model.Question, error) {
	questions, err := service.QuestionRepo.GetAll(context.Background())
	if err != nil {
		return make([]model.Question, 0), err
	}
	return questions, nil
}

func (service *QuestionService) Delete(questionID int) (bool, error) {
	deleted, err := service.QuestionRepo.Delete(context.Background(), questionID)
	if err != nil {
		return false, err
	}
	return deleted, nil
}
