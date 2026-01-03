package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	EnvTodoDBFile = "TODO_DBFILE"

	DefaultDBDir = "./data/"

	DefaultMigrationPath = "./migrations/scheduler.sql"
)

func Connect(dbFile string) (*sql.DB, error) {
	dbDir, err := getDBDirFromEnv()
	if err != nil {
		dbDir = DefaultDBDir
	}
	dbPath := filepath.Join(dbDir, dbFile)

	if err = os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dbDir, err)
	}

	_, err = os.Stat(dbPath)
	var needMigrateDB = false
	if err != nil {
		err = createDatabase(dbPath)
		if err != nil {
			return nil, err
		}
		needMigrateDB = true
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	
	if needMigrateDB {
		err = migrationDB(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func Disconnect(db *sql.DB) error {
	return db.Close()
}

func createDatabase(dbPath string) error {
	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}
	file.Close()

	return nil
}

func migrationDB(db *sql.DB) error {
	content, err := os.ReadFile(DefaultMigrationPath)
	if err != nil {
		return err
	}

	queries := strings.Split(string(content), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err = db.Exec(query)
		if err != nil {

			return fmt.Errorf("execute query %q: %w", query, err)
		}
	}

	return nil
}

func getDBDirFromEnv() (string, error) {
	todoDBFile, exist := os.LookupEnv(EnvTodoDBFile)
	if !exist {
		return "", fmt.Errorf("%s environment variable is not set", EnvTodoDBFile)
	}

	if todoDBFile == "" {
		return "", fmt.Errorf("%s is empty", EnvTodoDBFile)
	}

	return todoDBFile, nil
}
