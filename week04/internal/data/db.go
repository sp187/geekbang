package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/config"
	"sync"
)

var dbOnce sync.Once
var db *sql.DB

func InitDB(config *config.Config) (err error) {
	dbOnce.Do(func() {
		switch config.DbType {
		case "postgresql":
			dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
				config.DbUser,
				config.DbPassword,
				config.DbHost,
				config.DbPort,
				config.DbName,
			)
			db, err = sql.Open("postgres", dbURL)
			if err != nil {
				panic(fmt.Sprintf("connect to db %s fail: %s", dbURL, err.Error()))
			}
		//case "mysql":
		default:
			panic(fmt.Sprintf("unsupported db: %s", config.DbType))
		}
	})
	return nil
}

func GetDB() *sql.DB {
	if db == nil {
		panic("init db first")
	}
	return db
}
