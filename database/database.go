package database

import (
	"database/sql"
	"os"

	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDb() {
	var err error

	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		log.Fatal("Please set enviroment variable for DATABASE_URL")
	}

	db, err = sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()
}
