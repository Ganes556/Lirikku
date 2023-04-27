package models

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint `json:"id" gorm:"primarykey"` 
	CreatedAt *time.Time `json:",omitempty" gorm:"->:false;<-:create"`
	UpdatedAt *time.Time `json:",omitempty" gorm:"->:false;<-:create"`
	DeletedAt *gorm.DeletedAt `json:",omitempty" gorm:"->:false;<-:update"`
}