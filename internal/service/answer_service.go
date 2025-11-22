package service

import (
	"context"
	"database/sql"

	"github.com/mbilarusdev/quiz/internal/model"
	"github.com/mbilarusdev/quiz/internal/repository"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
	"gorm.io/gorm"
)

type AnswerLogic interface {
	AddAnswer(answer *model.Answer) (*model.Answer, error)
	FindAllByQuestionID(questionID int) ([]model.Answer, error)
	FindOne(answerID int) (*model.Answer, error)
	Delete(answerID int) (bool, error)
}

type AnswerService struct {
	AnswerRepo   repository.AnswerProvider
	QuestionRepo repository.QuestionProvider
	DB           *gorm.DB
}

func NewAnswerService(
	answerRepo repository.AnswerProvider,
	questionRepo repository.QuestionProvider,
	db *gorm.DB,
) *AnswerService {
	srv := new(AnswerService)
	srv.AnswerRepo = answerRepo
	srv.QuestionRepo = questionRepo
	srv.DB = db
	return srv
}

func (service *AnswerService) AddAnswer(answer *model.Answer) (*model.Answer, error) {
	ctx := context.Background()
	var newAnswer *model.Answer

	err := service.DB.Transaction(func(tx *gorm.DB) error {
		question, err := service.QuestionRepo.GetOne(tx, ctx, answer.QuestionID, false)
		if err != nil {
			return err
		}
		if question == nil {
			return &service_errors.NotFoundError{ID: answer.QuestionID}
		}

		newAnswer, err = service.AnswerRepo.Insert(tx, ctx, answer)
		if err != nil {
			return err
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	if err != nil {
		return nil, err
	}

	return newAnswer, nil
}

func (service *AnswerService) FindAllByQuestionID(questionID int) ([]model.Answer, error) {
	answers, err := service.AnswerRepo.GetAll(context.Background(), questionID)
	if err != nil {
		return make([]model.Answer, 0), err
	}
	return answers, nil
}

func (service *AnswerService) FindOne(answerID int) (*model.Answer, error) {
	answer, err := service.AnswerRepo.GetOne(context.Background(), answerID)
	if err != nil {
		return nil, err
	}
	if answer == nil {
		return answer, &service_errors.NotFoundError{ID: answerID}
	}
	return answer, nil
}

func (service *AnswerService) Delete(answerID int) (bool, error) {
	deleted, err := service.AnswerRepo.Delete(context.Background(), answerID)
	if err != nil {
		return false, err
	}
	return deleted, nil
}
