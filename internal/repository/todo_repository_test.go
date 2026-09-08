package repository

import (
	"path/filepath"
	"reflect"
	"testing"
	"todo/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdatePriorityPersistsAcrossDatabaseReopen(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "todo.db")
	db := openTestDatabase(t, dbPath)

	repository := NewTodoRepository(db)
	todo, err := repository.Create("persist priority")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated, err := repository.UpdatePriority(todo.ID, "high")
	if err != nil {
		t.Fatalf("UpdatePriority() error = %v", err)
	}
	if updated.Priority != "high" {
		t.Fatalf("UpdatePriority() priority = %q, want %q", updated.Priority, "high")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() error = %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reopened := openTestDatabase(t, dbPath)
	stored, err := NewTodoRepository(reopened).List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(stored) != 1 || stored[0].Priority != "high" {
		t.Fatalf("stored todos = %#v, want one todo with high priority", stored)
	}

	reopenedSQLDB, err := reopened.DB()
	if err != nil {
		t.Fatalf("reopened.DB() error = %v", err)
	}
	if err := reopenedSQLDB.Close(); err != nil {
		t.Fatalf("reopened Close() error = %v", err)
	}
}

func TestListUsesPriorityAsSortOrderTieBreaker(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "todo.db")
	db := openTestDatabase(t, dbPath)

	todos := []model.Todo{
		{ID: "low", NoteID: "1", Content: "low", Priority: "low", SortOrder: 1000, CreatedAt: "2026-01-01T00:00:00Z"},
		{ID: "high", NoteID: "1", Content: "high", Priority: "high", SortOrder: 1000, CreatedAt: "2026-01-01T00:00:00Z"},
		{ID: "later", NoteID: "1", Content: "later", Priority: "high", SortOrder: 2000, CreatedAt: "2026-01-01T00:00:00Z"},
	}
	if err := db.Create(&todos).Error; err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	listed, err := NewTodoRepository(db).List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	got := []string{listed[0].ID, listed[1].ID, listed[2].ID}
	want := []string{"high", "low", "later"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() IDs = %v, want %v", got, want)
	}
}

func openTestDatabase(t *testing.T, dbPath string) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}
	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
	return db
}
