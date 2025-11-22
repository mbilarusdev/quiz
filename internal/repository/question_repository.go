package repository

import (
	"context"
	"errors"

	"github.com/mbilarusdev/quiz/internal/model"
	service_errors "github.com/mbilarusdev/quiz/internal/service/errors"
	"github.com/mbilarusdev/quiz/internal/util"
	"gorm.io/gorm"
)

type QuestionProvider interface {
	Insert(ctx context.Context, question *model.Question) (*model.Question, error)
	GetOne(
		tx *gorm.DB,
		ctx context.Context,
		questionID int,
		withAnswers bool,
	) (*model.Question, error)
	GetAll(ctx context.Context) ([]model.Question, error)
	Delete(ctx context.Context, questionID int) (bool, error)
}

type QuestionRepository struct {
	DB *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	repo := new(QuestionRepository)
	repo.DB = db
	return repo
}

func (repo *QuestionRepository) Insert(
	ctx context.Context,
	question *model.Question,
) (*model.Question, error) {
	if err := gorm.G[model.Question](repo.DB).Create(ctx, question); err != nil {
		if util.CheckDublicateErr(err) {
			return nil, &service_errors.DuplicateError{ID: question.ID}
		}
		return nil, err
	}
	return question, nil
}

func (repo *QuestionRepository) GetOne(
	tx *gorm.DB,
	ctx context.Context,
	questionID int,
	withAnswers bool,
) (*model.Question, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.DB
	}
	iface := gorm.G[model.Question](db)

	query := iface.Where("id = ?", questionID)
	if withAnswers {
		query = query.Preload("Answers", func(_ gorm.PreloadBuilder) error {
			return db.Order("created_at ASC").Error
		})
	}

	question, err := query.First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

func (repo *QuestionRepository) GetAll(ctx context.Context) ([]model.Question, error) {
	questions, err := gorm.G[model.Question](
		repo.DB,
	).Find(ctx)
	if err != nil {
		return make([]model.Question, 0), nil
	}
	return questions, nil
}

func (repo *QuestionRepository) Delete(ctx context.Context, questionID int) (bool, error) {
	rowsAffected, err := gorm.G[model.Question](repo.DB).Where("id = ?", questionID).Delete(ctx)
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
