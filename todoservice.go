package main

import (
	"todo/internal/model"
	"todo/internal/repository"
)

type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) List() ([]model.Todo, error) {
	return s.repo.List()
}

func (s *TodoService) Create(content string) (*model.Todo, error) {
	return s.repo.Create(content)
}

func (s *TodoService) Toggle(id string, completed bool) error {
	return s.repo.SetCompleted(id, completed)
}

func (s *TodoService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *TodoService) UpdateContent(id string, content string) error {
	return s.repo.UpdateContent(id, content)
}
