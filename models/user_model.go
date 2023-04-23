package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) CheckPassword(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
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