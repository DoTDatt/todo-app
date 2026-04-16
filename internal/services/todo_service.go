package services

import (
	"errors"

	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/repositories"
)

type TodoService struct {
	repo *repositories.TodoRepository
}

func NewtodoService(repo *repositories.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("Title không được trống")
	}
	validStatuses := map[string]bool{"pending": true, " inProgress ": true, " done": true}
	if !validStatuses[todo.Status] {
		return errors.New("status không hợp lệ(chỉ chấp nhận : pending , inProgress , done)")
	}

	return s.repo.Create(todo)
}

func (s *TodoService) GetbyID(id uint) (*models.Todo, error) {
	return s.repo.GetbyID(id)
}

func (s *TodoService) GetAll() ([]models.Todo, error) {
	return s.repo.GetAll()
}

func (s *TodoService) Update(todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("Title không được trống")
	}
	validStatuses := map[string]bool{"pending": true, " inProgress ": true, " done": true}
	if !validStatuses[todo.Status] {
		return errors.New("status không hợp lệ(chỉ chấp nhận : pending , inProgress , done)")
	}

	return s.repo.Update(todo)
}

func (s *TodoService) Delete(id uint) error {
	return s.repo.Delete(id)
}
