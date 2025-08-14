package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetDB(databasePath string) *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", databasePath))
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
			os.Exit(1)
		}
		if err = db.Ping(); err != nil {
			fmt.Printf("Failed to ping database: %v\n", err)
			os.Exit(1)
		}
	})
	return db
}
