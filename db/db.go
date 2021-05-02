package db

import (
	"database/sql"
	"log"
	"time"
)

var db *sql.DB

func GetConnection() *sql.DB {
	if db == nil {
		// log.Println("database connection lost, reconenct...")
		conn, err := sql.Open("sqlite3", "./app.db")
		if err != nil {
			log.Fatal("Open database error", err.Error())
		}
		db = conn
	} else {
		// log.Println("database connection maintained")
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

func CloseConnection() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatal("Close database error", err.Error())
		}
	}
}
