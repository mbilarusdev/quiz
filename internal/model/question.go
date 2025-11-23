package model

import (
	"time"

	"go.uber.org/zap/zapcore"
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

func (q Question) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", q.ID)
	enc.AddString("text", q.Text)
	enc.AddTime("created_at", q.CreatedAt)
	enc.AddTime("updated_at", q.UpdatedAt)
	enc.AddTime("deleted_at", q.DeletedAt.Time)
	enc.AddArray("answers", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, ans := range q.Answers {
			enc.AppendObject(ans)
		}
		return nil
	}))
	return nil
}
