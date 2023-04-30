package configs

import (
	"fmt"
	"os"

	"github.com/Lirikku/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
	os.Getenv("DB_USERNAME"), 
	os.Getenv("DB_PASSWORD"), 
	os.Getenv("DB_HOST"), 
	os.Getenv("DB_NAME"))
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	initMigrate()
}

func initMigrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.SongLyric{})
}