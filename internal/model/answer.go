package model

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type Answer struct {
	ID         int            `gorm:"primarykey,autoIncrement:true" json:"id,omitempty"`
	QuestionID int            `gorm:"index;not null"                json:"question_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null"      json:"user_id"`
	Text       string         `gorm:"not null"                      json:"text"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"                json:"created_at,omitzero"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"                json:"updated_at,omitzero"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (a Answer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", a.ID)
	enc.AddInt("question_id", a.QuestionID)
	enc.AddString("text", a.Text)
	enc.AddTime("created_at", a.CreatedAt)
	enc.AddTime("updated_at", a.UpdatedAt)
	enc.AddTime("deleted_at", a.DeletedAt.Time)
	return nil
}
