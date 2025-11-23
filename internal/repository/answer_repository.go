package repository

import (
	"context"
	"errors"

	"github.com/mbilarusdev/quiz/internal/common"
	"github.com/mbilarusdev/quiz/internal/model"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
	"github.com/mbilarusdev/quiz/internal/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AnswerProvider interface {
	Insert(
		tx *gorm.DB,
		ctx context.Context,
		answer *model.Answer,
	) (*model.Answer, error)
	GetOne(ctx context.Context, answerID int) (*model.Answer, error)
	Delete(ctx context.Context, answerID int) (bool, error)
}

type AnswerRepository struct {
	DB *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	repo := new(AnswerRepository)
	repo.DB = db
	return repo
}

func (repo *AnswerRepository) Insert(
	tx *gorm.DB,
	ctx context.Context,
	answer *model.Answer,
) (*model.Answer, error) {
	op := "repository.AnswerRepository.Insert"
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.DB
	}
	if err := gorm.G[model.Answer](db).Create(ctx, answer); err != nil {
		if util.CheckDublicateErr(err) {
			common.L.Error("DB error",
				zap.String("op", op),
				zap.String("Result", "Duplicated key when create answer"),
				zap.Object("Answer", answer))
			return nil, &service_errors.DuplicateError{ID: answer.ID}
		}
		common.L.Error("DB error",
			zap.String("op", op),
			zap.String("Result", "Error occured when create answer"),
			zap.Object("Answer", answer))
		return nil, err
	}
	common.L.Info("DB success",
		zap.String("op", op),
		zap.String("Result", "Answer created successfully!"),
		zap.Object("Answer", answer))
	return answer, nil
}

func (repo *AnswerRepository) GetOne(ctx context.Context, answerID int) (*model.Answer, error) {
	op := "repository.AnswerRepository.GetOne"
	answer, err := gorm.G[model.Answer](repo.DB).Where("id = ?", answerID).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.L.Warn("DB warn",
				zap.String("op", op),
				zap.String("Result", "Answer not found"),
				zap.Int("answer_id", answerID))
			return nil, nil
		}
		common.L.Error("DB error",
			zap.String("op", op),
			zap.String("Result", "Error occured when find one answer"),
			zap.Int("answer_id", answerID))
		return nil, err
	}
	common.L.Info("DB success",
		zap.String("op", op),
		zap.String("Result", "Answer finded successfully!"),
		zap.Object("Answer", answer))
	return &answer, nil
}

func (repo *AnswerRepository) Delete(ctx context.Context, answerID int) (bool, error) {
	op := "repository.AnswerRepository.Delete"
	rowsAffected, err := gorm.G[model.Answer](repo.DB).Where("id = ?", answerID).Delete(ctx)
	if err != nil {
		common.L.Error("DB error",
			zap.String("op", op),
			zap.String("Result", "Error occured when delete Answer"),
			zap.Int("answer_id", answerID))
		return false, err
	}
	if rowsAffected == 0 {
		common.L.Warn("DB warn",
			zap.String("op", op),
			zap.String("Result", "Deletable answer not found"),
			zap.Int("answer_id", answerID))
		return false, nil
	}
	common.L.Info("DB info",
		zap.String("op", op),
		zap.String("Result", "Answer deleted with success"),
		zap.Int("answer_id", answerID))
	return true, nil
}
