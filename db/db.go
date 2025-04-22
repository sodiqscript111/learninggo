package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // ✅ CGO-free SQLite driver
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "C:/Users/Owner/GolandProjects/learninggo/api.db")

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTable()
	log.Println("✅ Database initialized and table ready.")
}

func createTable() {
	createEventTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INTEGER
	);`

	_, err := DB.Exec(createEventTable)
	if err != nil {
		log.Fatalf("❌ Failed to create events table: %v", err)
	}
}
