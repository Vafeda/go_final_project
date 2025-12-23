package database

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"strings"
)

func Connect(dbFile string) (*sql.DB, error) {
	dbFile = "../../data/" + dbFile
	_, err := os.Stat(dbFile)
	var migrateDB = false
	if err != nil {
		dbFile, err = createDatabase()
		if err != nil {
			return nil, err
		}
		migrateDB = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	if migrateDB {
		migrationDB(db)
	}

	return db, nil
}

func Disconnect(db *sql.DB) error {
	return db.Close()
}

func createDatabase() (string, error) {
	_, err := os.Create("data/scheduler.db")
	if err != nil {
		return "", err
	}

	return "data/scheduler.db", nil
}

func migrationDB(db *sql.DB) error {
	content, err := os.ReadFile("./migrations/scheduler.sql")
	if err != nil {
		return err
	}

	queries := strings.Split(string(content), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("execute query %q: %w", query, err)
		}
	}

	return nil
}
