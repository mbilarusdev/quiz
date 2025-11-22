package model

import (
	"time"

	"github.com/google/uuid"
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
