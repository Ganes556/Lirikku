package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID      uint `json:"id" gorm:"primarykey"` 
	Name 	  string `json:"name" gorm:"type:varchar(150)"`
	Email 	 string `json:"email" gorm:"type:varchar(255)"`
	Password string `json:"password" gorm:"type:varchar(64)"`
	SongLyrics 	[]SongLyric `json:"song_lyrics" gorm:"foreignKey:UserID"`
	Base
}

type UserRegister struct {
	Name 	 string `json:"name"`
	Email 	string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email 	string `json:"email"`
	Password string `json:"password"`
}

type UserJWTDecode struct {
	ID uint
	Name string
}

type JWTClaims struct {
	ID    uint   `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hash)
	return
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

func (u *User) GenerateToken() (string, error) {
	claims := &JWTClaims{
		ID: u.ID,
		Name: u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}