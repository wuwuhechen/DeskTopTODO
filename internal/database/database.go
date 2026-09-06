package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Open() (*sql.DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %v", err)
	}

	appDir := filepath.Join(configDir, "DesktopTodo")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app directory: %v", err)
	}

	dbPath := filepath.Join(appDir, "todo.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		return nil, err
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		sort_order INTEGER NOT NULL,
		created_at TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_todos_sort_order ON todos (sort_order);
	`)
	return err
}
