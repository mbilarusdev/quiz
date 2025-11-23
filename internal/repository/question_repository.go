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
	op := "repository.QuestionRepository.Insert"
	if err := gorm.G[model.Question](repo.DB).Create(ctx, question); err != nil {
		if util.CheckDublicateErr(err) {
			common.L.Error("DB error",
				zap.String("op", op),
				zap.String("Result", "Duplicated key when create question"),
				zap.Object("Question", question))
			return nil, &service_errors.DuplicateError{ID: question.ID}
		}
		common.L.Error("DB error",
			zap.String("op", op),
			zap.String("Result", "Error occured when create question"),
			zap.Object("Question", question))
		return nil, err
	}
	common.L.Info("DB success",
		zap.String("op", op),
		zap.String("Result", "Question created successfully!"),
		zap.Object("Question", question))
	return question, nil
}

func (repo *QuestionRepository) GetOne(
	tx *gorm.DB,
	ctx context.Context,
	questionID int,
	withAnswers bool,
) (*model.Question, error) {
	op := "repository.QuestionRepository.GetOne"
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
			common.L.Warn("DB warn",
				zap.String("op", op),
				zap.String("Result", "Question not found"),
				zap.Int("question_id", questionID))
			return nil, nil
		}
		common.L.Error("DB error",
			zap.String("op", op),
			zap.String("Result", "Error occured when find one question with answers"),
			zap.Int("question_id", questionID))
		return nil, err
	}
	common.L.Info("DB success",
		zap.String("op", op),
		zap.String("Result", "Question with answers finded successfully!"),
		zap.Object("Question", question))
	return &question, nil
}

func (repo *QuestionRepository) GetAll(ctx context.Context) ([]model.Question, error) {
	op := "repository.QuestionRepository.GetAll"
	questions, err := gorm.G[model.Question](
		repo.DB,
	).Find(ctx)
	if err != nil {
		common.L.Error("DB error",
			zap.String("op", op),
			zap.String("Result", "Error occured when try to find all questions"))
		return make([]model.Question, 0), nil
	}
	common.L.Error("DB success",
		zap.String("op", op),
		zap.String("Result", "All questions finded with success"))
	return questions, nil
}

func (repo *QuestionRepository) Delete(ctx context.Context, questionID int) (bool, error) {
	op := "repository.QuestionRepository.Delete"
	rowsAffected, err := gorm.G[model.Question](repo.DB).Where("id = ?", questionID).Delete(ctx)
	if err != nil {
		common.L.Error("DB error",
			zap.String("op", op),
			zap.String("Result", "Error occured when delete Question"),
			zap.Int("question_id", questionID))
		return false, err
	}
	if rowsAffected == 0 {
		common.L.Warn("DB warn",
			zap.String("op", op),
			zap.String("Result", "Deletable question not found"),
			zap.Int("question_id", questionID))
		return false, nil
	}
	common.L.Info("DB info",
		zap.String("op", op),
		zap.String("Result", "Question deleted with success"),
		zap.Int("question_id", questionID))
	return true, nil
}
