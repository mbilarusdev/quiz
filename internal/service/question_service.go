package service

import (
	"context"

	"github.com/mbilarusdev/quiz/internal/common"
	"github.com/mbilarusdev/quiz/internal/model"
	"github.com/mbilarusdev/quiz/internal/repository"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
	"go.uber.org/zap"
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
	op := "service.QuestionService.Create"
	question, err := service.QuestionRepo.Insert(context.Background(), question)
	if err != nil {
		common.L.Error("Domain error",
			zap.String("op", op),
			zap.String("Result", "Error when try to create question"))
		return nil, err
	}
	common.L.Info("Domain info",
		zap.String("op", op),
		zap.String("Result", "Question created with success"),
		zap.Int("id", question.ID))
	return question, nil
}

func (service *QuestionService) FindOneDetailed(questionID int) (*model.Question, error) {
	op := "service.QuestionService.FindOneDetailed"
	ctx := context.Background()
	question, err := service.QuestionRepo.GetOne(nil, ctx, questionID, true)
	if err != nil {
		common.L.Error("Domain error",
			zap.String("op", op),
			zap.String("Result", "Error when try to find question"),
			zap.Int("id", questionID))
		return nil, err
	}
	if question == nil {
		common.L.Warn("Domain warn",
			zap.String("op", op),
			zap.String("Result", "Question not found!"),
			zap.Int("id", questionID))
		return question, &service_errors.NotFoundError{ID: questionID}
	}

	common.L.Info("Domain info",
		zap.String("op", op),
		zap.String("Result", "Question finded with success"),
		zap.Int("id", questionID))
	return question, nil
}

func (service *QuestionService) FindAll() ([]model.Question, error) {
	op := "service.QuestionService.FindOneDetailed"
	questions, err := service.QuestionRepo.GetAll(context.Background())
	if err != nil {
		common.L.Error("Domain error",
			zap.String("op", op),
			zap.String("Result", "Error when try to find all of questions"))
		return make([]model.Question, 0), err
	}
	common.L.Info("Domain info",
		zap.String("op", op),
		zap.String("Result", "All of questions finded successfully"))
	return questions, nil
}

func (service *QuestionService) Delete(questionID int) (bool, error) {
	op := "service.QuestionService.Delete"
	deleted, err := service.QuestionRepo.Delete(context.Background(), questionID)
	if err != nil {
		common.L.Error("Domain error",
			zap.String("op", op),
			zap.String("Result", "Error when try to delete question"),
			zap.Int("id", questionID))
		return false, err
	}
	if !deleted {
		common.L.Warn("Domain warn",
			zap.String("op", op),
			zap.String("Result", "Question to delete not found!"),
			zap.Int("id", questionID))
	}
	common.L.Info("Domain info",
		zap.String("op", op),
		zap.String("Result", "Question deleted successfully!"),
		zap.Int("id", questionID))
	return deleted, nil
}
