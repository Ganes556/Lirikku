package configs

import (
	"os"
	"time"

	"github.com/srinathgs/mysqlstore"
)

var Store *mysqlstore.MySQLStore

func InitSessionMysql() {
	if SqlCon == nil {
		panic("please init db con first!")
	}
	var err error
	Store, err = mysqlstore.NewMySQLStoreFromConnection(
		SqlCon,
		"session",
		"/",	
		int(time.Duration(24*time.Hour*7).Seconds()), []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		panic(err)
	}
}