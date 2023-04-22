package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name 	 string `json:"name" gorm:"type:varchar(150)"`
	Email 	string `json:"email" gorm:"type:varchar(255)"`
	Password string `json:"password" gorm:"type:varchar(64)"`
	Songs 	[]SongLyric `json:"songs" gorm:"foreignKey:UserID"`
}
