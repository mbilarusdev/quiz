package service

import (
	"context"
	"database/sql"

	"github.com/mbilarusdev/quiz/internal/common"
	"github.com/mbilarusdev/quiz/internal/model"
	"github.com/mbilarusdev/quiz/internal/repository"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AnswerLogic interface {
	AddAnswer(answer *model.Answer) (*model.Answer, error)
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
	op := "service.AnswerService.AddAnswer"
	ctx := context.Background()
	var newAnswer *model.Answer

	err := service.DB.Transaction(func(tx *gorm.DB) error {
		question, err := service.QuestionRepo.GetOne(tx, ctx, answer.QuestionID, false)
		if err != nil {
			common.L.Error("Domain error",
				zap.String("op", op),
				zap.String("Result", "Error when try to find question which needed add answer"))
			return err
		}
		if question == nil {
			common.L.Warn("Domain warn",
				zap.String("op", op),
				zap.String("Result", "Question which needed to add answer not found"))
			return &service_errors.NotFoundError{ID: answer.QuestionID}
		}

		newAnswer, err = service.AnswerRepo.Insert(tx, ctx, answer)
		if err != nil {
			common.L.Error("Domain error",
				zap.String("op", op),
				zap.String("Result", "Error occured when inserting answer"))
			return err
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	if err != nil {
		return nil, err
	}

	common.L.Info("Domain info",
		zap.String("op", op),
		zap.String("Result", "Answer added successfully"),
		zap.Int("id", newAnswer.ID),
		zap.Int("question_id", newAnswer.QuestionID))
	return newAnswer, nil
}

func (service *AnswerService) FindOne(answerID int) (*model.Answer, error) {
	op := "service.AnswerService.AddAnswer"
	answer, err := service.AnswerRepo.GetOne(context.Background(), answerID)
	if err != nil {
		common.L.Info("Domain error",
			zap.String("op", op),
			zap.String("Result", "Error when try to find answer"),
			zap.Int("id", answerID))
		return nil, err
	}
	if answer == nil {
		common.L.Warn("Domain warn",
			zap.String("op", op),
			zap.String("Result", "Answer not found"),
			zap.Int("id", answerID))
		return answer, &service_errors.NotFoundError{ID: answerID}
	}
	common.L.Info("Domain info",
		zap.String("op", op),
		zap.String("Result", "Answer finded successfully!"),
		zap.Int("id", answerID))
	return answer, nil
}

func (service *AnswerService) Delete(answerID int) (bool, error) {
	op := "service.AnswerService.Delete"
	deleted, err := service.AnswerRepo.Delete(context.Background(), answerID)
	if err != nil {
		common.L.Error("Domain error",
			zap.String("op", op),
			zap.String("Result", "Error when try to delete answer"),
			zap.Int("id", answerID))
		return false, err
	}
	if !deleted {
		common.L.Warn("Domain warn",
			zap.String("op", op),
			zap.String("Result", "Answer to delete not found!"),
			zap.Int("id", answerID))
	}
	common.L.Info("Domain info",
		zap.String("op", op),
		zap.String("Result", "Answer deleted successfully!"),
		zap.Int("id", answerID))
	return deleted, nil
}
