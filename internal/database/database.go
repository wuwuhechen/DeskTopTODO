package database

import (
	"fmt"
	"os"
	"path/filepath"
	"todo/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open() (*gorm.DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	appDir := filepath.Join(configDir, "DesktopTodo")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app directory: %w", err)
	}

	dbPath := filepath.Join(appDir, "todo.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	defaultSetting := model.Setting{Key: "show_completed", Value: "false"}
	if err := db.Where("key = ?", defaultSetting.Key).FirstOrCreate(&defaultSetting).Error; err != nil {
		return nil, fmt.Errorf("failed to initialise settings: %w", err)
	}
	return db, nil
}

func migrate(db *gorm.DB) error {
	migrator := db.Migrator()

	if !migrator.HasTable(&model.Todo{}) {
		if err := migrator.CreateTable(&model.Todo{}); err != nil {
			return err
		}
	} else {
		for _, field := range []string{"NoteID", "Priority", "UpdatedAt"} {
			if !migrator.HasColumn(&model.Todo{}, field) {
				if err := migrator.AddColumn(&model.Todo{}, field); err != nil {
					return err
				}
			}
		}
	}

	if !migrator.HasIndex(&model.Todo{}, "idx_todos_sort_order") {
		if err := migrator.CreateIndex(&model.Todo{}, "SortOrder"); err != nil {
			return err
		}
	}

	if !migrator.HasTable(&model.Setting{}) {
		if err := migrator.CreateTable(&model.Setting{}); err != nil {
			return err
		}
	}
	return nil
}
