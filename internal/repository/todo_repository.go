package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"todo/internal/model"
	"unicode/utf8"

	"github.com/google/uuid"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) List() ([]model.Todo, error) {
	rows, err := r.db.Query(`
	     SELECT id, content, completed, priority, sort_order, created_at, updated_at
		 FROM todos
		 ORDER BY sort_order ASC, created_at ASC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]model.Todo, 0)

	for rows.Next() {
		var todo model.Todo
		var completed bool

		if err := rows.Scan(
			&todo.ID,
			&todo.Content,
			&completed,
			&todo.Priority,
			&todo.SortOrder,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, err
		}

		todo.Completed = completed
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (r *TodoRepository) Create(content string) (*model.Todo, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return &model.Todo{}, fmt.Errorf("content cannot be empty")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return &model.Todo{}, err
	}
	defer tx.Rollback()

	var nextOrder int
	if err := tx.QueryRow(`
		SELECT COALESCE(MAX(sort_order), -1) + 1 FROM todos
	`).Scan(&nextOrder); err != nil {
		return &model.Todo{}, err
	} // TODO

	todo := &model.Todo{
		ID:        uuid.New().String(),
		Content:   content,
		Completed: false,
		Priority:  "normal",
		SortOrder: nextOrder,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: "",
	}

	_, err = tx.Exec(`
	    INSERT INTO todos (id, note_id, content, completed, priority, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, todo.ID, todo.NoteID, todo.Content, 0, todo.Priority, todo.SortOrder, todo.CreatedAt, todo.UpdatedAt)

	if err != nil {
		return &model.Todo{}, err
	}

	if err := tx.Commit(); err != nil {
		return &model.Todo{}, err
	}

	return todo, nil
}

func (r *TodoRepository) SetCompleted(id string, completed bool) error {
	result, err := r.db.Exec(`
		UPDATE todos
		SET completed = ?
		WHERE id = ?
	`, completed, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("todo with id %s not found", id)
	}
	return nil
}

func (r *TodoRepository) Delete(id string) error {
	result, err := r.db.Exec(`
		DELETE FROM todos
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("todo with id %s not found", id)
	}
	return nil
}

func (r *TodoRepository) UpdateContent(id string, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	if utf8.RuneCountInString(content) > 500 {
		return fmt.Errorf("content cannot exceed 500 characters")
	}

	_, err := r.db.Exec(`
		UPDATE todos
		SET content = ?, updated_at = ?
		WHERE id = ?
	`, content, time.Now().UTC().Format(time.RFC3339), id)

	return err
}

func (r *TodoRepository) SetShowCompleted(showCompleted string) error {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM settings WHERE key = 'show_completed'
	`).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		_, err = r.db.Exec(`
			INSERT INTO settings (key, value) VALUES ('show_completed', ?)
		`, showCompleted)
	} else {
		_, err = r.db.Exec(`
			UPDATE settings SET value = ? WHERE key = 'show_completed'
		`, showCompleted)
	}

	return err
}

func (r *TodoRepository) GetShowCompleted() (string, error) {
	var value string
	err := r.db.QueryRow(`
		SELECT value FROM settings WHERE key = 'show_completed'
	`).Scan(&value)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *TodoRepository) UpdatePriority(id string, priority string) error {
	switch priority {
	case "normal", "low", "medium", "high":
	default:
		return fmt.Errorf("invalid priority value: %s", priority)
	}

	_, err := r.db.Exec(`
		UPDATE todos
		SET priority = ?, updated_at = ?
		WHERE id = ?
	`, priority, time.Now().UTC().Format(time.RFC3339), id)

	return err
}

func (r *TodoRepository) Recorder(noteID string, ids []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for idx, id := range ids {
		result, err := tx.Exec(`
			UPDATE todos
			SET sort_order = ?, updated_at = ?
			WHERE id = ? AND note_id = ?
		`, (idx+1)*1000, time.Now().UTC().Format(time.RFC3339), id, noteID)
		if err != nil {
			return err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			return fmt.Errorf("todo with id %s not found for note_id %s", id, noteID)
		}
	}
	return tx.Commit()
}
