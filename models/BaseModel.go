package models

import (
	"time"
)

type BaseModel struct {
	Id int `json:"id" gorm:"AUTO_INCREMENT; primary_key"`
	CreatedAt time.Time `json:"created_at" sql:"index"`
	UpdatedAt time.Time `json:"updated_at" sql:"index"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}