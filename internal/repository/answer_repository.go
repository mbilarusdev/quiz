package repository

import (
	"context"
	"errors"

	"github.com/mbilarusdev/quiz/internal/model"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
	"github.com/mbilarusdev/quiz/internal/util"
	"gorm.io/gorm"
)

type AnswerProvider interface {
	Insert(
		tx *gorm.DB,
		ctx context.Context,
		answer *model.Answer,
	) (*model.Answer, error)
	GetOne(ctx context.Context, answerID int) (*model.Answer, error)
	GetAll(ctx context.Context, questionID int) ([]model.Answer, error)
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
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.DB
	}
	if err := gorm.G[model.Answer](db).Create(ctx, answer); err != nil {
		if util.CheckDublicateErr(err) {
			return nil, &service_errors.DuplicateError{ID: answer.ID}
		}
		return nil, err
	}
	return answer, nil
}

func (repo *AnswerRepository) GetOne(ctx context.Context, answerID int) (*model.Answer, error) {
	answer, err := gorm.G[model.Answer](repo.DB).Where("id = ?", answerID).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &answer, nil
}

func (repo *AnswerRepository) GetAll(
	ctx context.Context,
	questionID int,
) ([]model.Answer, error) {
	query := gorm.G[model.Answer](
		repo.DB,
	).Where("TRUE")

	if questionID != 0 {
		query = query.Where("question_id = ?", questionID)
	}

	answers, err := query.Find(ctx)
	if err != nil {
		return make([]model.Answer, 0), nil
	}
	return answers, nil
}

func (repo *AnswerRepository) Delete(ctx context.Context, answerID int) (bool, error) {
	rowsAffected, err := gorm.G[model.Answer](repo.DB).Where("id = ?", answerID).Delete(ctx)
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
