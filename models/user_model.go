package models

import (
	"github.com/Lirikku/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Name 	 string `json:"name" gorm:"type:varchar(150)"`
	Email 	string `json:"email" gorm:"type:varchar(255)"`
	Password string `json:"password" gorm:"type:varchar(64)"`
	Songs 	[]SongLyric `json:"songs" gorm:"foreignKey:UserID"`
}

type ReqAuthUser struct {
	Name 	 string `json:"name"`
	Email 	string `json:"email"`
	Password string `json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password = utils.HashPassword(u.Password)
	return
}