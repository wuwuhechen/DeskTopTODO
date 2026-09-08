package repository

import (
	"fmt"
	"log"
	"strings"
	"time"
	"todo/internal/model"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) List() ([]model.Todo, error) {
	var todos []model.Todo
	err := r.db.
		Order(`CASE priority
			WHEN 'high' THEN 1
			WHEN 'medium' THEN 2
			WHEN 'low' THEN 3
			WHEN 'normal' THEN 4
			ELSE 5
		END ASC`).
		Order("sort_order ASC").
		Order("created_at ASC").
		Find(&todos).Error
	return todos, err
}

func (r *TodoRepository) Create(content string) (*model.Todo, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	todo := &model.Todo{
		ID: uuid.NewString(), NoteID: "1", Content: content, Priority: "normal",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var maxOrder int
		if err := tx.Model(&model.Todo{}).Select("COALESCE(MAX(sort_order), -1)").Scan(&maxOrder).Error; err != nil {
			return err
		}
		todo.SortOrder = maxOrder + 1
		return tx.Create(todo).Error
	})
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *TodoRepository) SetCompleted(id string, completed bool) error {
	return r.requireTodo(r.db.Model(&model.Todo{}).Where("id = ?", id).Update("completed", completed), id)
}

func (r *TodoRepository) Delete(id string) error {
	return r.requireTodo(r.db.Delete(&model.Todo{}, "id = ?", id), id)
}

func (r *TodoRepository) UpdateContent(id, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if utf8.RuneCountInString(content) > 500 {
		return fmt.Errorf("content cannot exceed 500 characters")
	}
	return r.requireTodo(r.db.Model(&model.Todo{}).Where("id = ?", id).Updates(map[string]any{
		"content": content, "updated_at": time.Now().UTC().Format(time.RFC3339),
	}), id)
}

func (r *TodoRepository) UpdatePriority(id, priority string) (*model.Todo, error) {
	switch priority {
	case "normal", "low", "medium", "high":
	default:
		return nil, fmt.Errorf("invalid priority value: %s", priority)
	}

	log.Printf("Updating priority for todo ID %s to %s", id, priority)

	result := r.db.Model(&model.Todo{}).Where("id = ?", id).Updates(map[string]any{
		"priority": priority, "updated_at": time.Now().UTC().Format(time.RFC3339),
	})
	if err := r.requireTodo(result, id); err != nil {
		log.Printf("Priority update failed for todo ID %s: %v", id, err)
		return nil, err
	}

	var todo model.Todo
	if err := r.db.First(&todo, "id = ?", id).Error; err != nil {
		log.Printf("Priority update read-back failed for todo ID %s: %v", id, err)
		return nil, err
	}
	log.Printf("Priority update persisted for todo ID %s: priority=%s updated_at=%s", todo.ID, todo.Priority, todo.UpdatedAt)
	return &todo, nil
}

func (r *TodoRepository) Recorder(noteID string, ids []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for idx, id := range ids {
			result := tx.Model(&model.Todo{}).Where("id = ? AND note_id = ?", id, noteID).Updates(map[string]any{
				"sort_order": (idx + 1) * 1000,
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("todo with id %s not found for note_id %s", id, noteID)
			}
		}
		return nil
	})
}

func (r *TodoRepository) requireTodo(result *gorm.DB, id string) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("todo with id %s not found", id)
	}
	return nil
}
