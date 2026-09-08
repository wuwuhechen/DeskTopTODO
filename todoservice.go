package main

import (
	"todo/internal/model"
	"todo/internal/repository"
)

type TodoService struct {
	todos    *repository.TodoRepository
	settings *repository.SettingRepository
}

func NewTodoService(todos *repository.TodoRepository, settings *repository.SettingRepository) *TodoService {
	return &TodoService{todos: todos, settings: settings}
}

func (s *TodoService) List() ([]model.Todo, error) {
	return s.todos.List()
}

func (s *TodoService) Create(content string) (*model.Todo, error) {
	return s.todos.Create(content)
}

func (s *TodoService) Toggle(id string, completed bool) error {
	return s.todos.SetCompleted(id, completed)
}

func (s *TodoService) Delete(id string) error {
	return s.todos.Delete(id)
}

func (s *TodoService) UpdateContent(id string, content string) error {
	return s.todos.UpdateContent(id, content)
}

func (s *TodoService) SetShowCompleted(showCompleted string) error {
	return s.settings.SetShowCompleted(showCompleted)
}

func (s *TodoService) GetShowCompleted() (string, error) {
	return s.settings.GetShowCompleted()
}

func (s *TodoService) UpdatePriority(id string, priority string) (*model.Todo, error) {
	return s.todos.UpdatePriority(id, priority)
}

func (s *TodoService) UpdateSortOrder(noteID string, ids []string) error {
	return s.todos.Recorder(noteID, ids)
}
