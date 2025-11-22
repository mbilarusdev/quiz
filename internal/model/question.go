package model

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID        int            `gorm:"primarykey"                                                          json:"id,omitempty"`
	Text      string         `gorm:"not null"                                                            json:"text"`
	CreatedAt time.Time      `gorm:"autoCreateTime"                                                      json:"created_at,omitzero"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"                                                      json:"updated_at,omitzero"`
	Answers   []Answer       `gorm:"foreignkey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"answers,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
