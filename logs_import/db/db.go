package db

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/abaldeweg/warehouse-server/logs/entity"
	_ "github.com/mattn/go-sqlite3"
)

// DBHandler handles database operations for logs.
type DBHandler struct {
	db *sql.DB
}

// NewDBHandler creates a new DBHandler.
func NewDBHandler() (*DBHandler, error) {
	setup()

	db, err := sql.Open("sqlite3", "data/db/events.db")
	if err != nil {
		return nil, err
	}

	migrate(db)

	return &DBHandler{db: db}, nil
}

// Close closes the database connection.
func (handler *DBHandler) Close() error {
	return handler.db.Close()
}

// Write inserts a log entry into the database.
func (handler *DBHandler) Write(date int, data entity.LogEntry) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	query := `INSERT INTO logs (date, data) VALUES (?, ?)`
	_, err = handler.db.Exec(query, date, jsonData)
	return err
}

// Exists checks if a log entry already exists in the database.
func (handler *DBHandler) Exists(date int, data entity.LogEntry) (bool, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	query := `SELECT COUNT(*) FROM logs WHERE date = ? AND data = ?`
	var count int
	err = handler.db.QueryRow(query, date, jsonData).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// setup creates the necessary directories and files for the database.
func setup() error {
	if err := os.MkdirAll("data/db", os.ModePerm); err != nil {
		return err
	}
	_, err := os.Stat("data/db/events.db")
	if os.IsNotExist(err) {
		file, err := os.Create("data/db/events.db")
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

// migrate creates the necessary tables for the database.
func migrate(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date INTEGER,
		data TEXT
	)`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}
