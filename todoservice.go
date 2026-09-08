package main

import (
	"todo/internal/model"
	"todo/internal/repository"
)

type TodoService struct {
	todos *repository.TodoRepository
}

func NewTodoService(todos *repository.TodoRepository) *TodoService {
	return &TodoService{todos: todos}
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

func (s *TodoService) UpdatePriority(id string, priority string) (*model.Todo, error) {
	return s.todos.UpdatePriority(id, priority)
}

func (s *TodoService) UpdateSortOrder(noteID string, ids []string) error {
	return s.todos.Recorder(noteID, ids)
}
