package configs

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Lirikku/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SqlCon *sql.DB
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
	os.Getenv("DB_USERNAME"), 
	os.Getenv("DB_PASSWORD"), 
	os.Getenv("DB_HOST"), 
	os.Getenv("DB_NAME"))
	var err error
	SqlCon, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: SqlCon,
	}))
	if err != nil {
		panic(err)
	}
	initMigrate()
}

func initMigrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.SongLyric{})
}