package services

import (
	"errors"

	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/repositories"
	"gorm.io/gorm"
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
	validStatuses := map[string]bool{"pending": true, "in_progress": true, "done": true}
	if !validStatuses[todo.Status] {
		return errors.New("status không hợp lệ(chỉ chấp nhận : pending , in_progress , done)")
	}

	return s.repo.Create(todo)
}

func (s *TodoService) GetbyID(id uint) (*models.Todo, error) {
	todo, err := s.repo.GetbyID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Todo không tồn tại")
	}
	return todo, nil
}

func (s *TodoService) GetAll(userID int, status string, search string, sort string, order string) ([]models.Todo, error) {
	if sort == "" {
		sort = "id"
	}

	if order == "" {
		order = "desc"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	return s.repo.GetAll(userID, status, search, sort, order)
}

func (s *TodoService) Update(todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("Title không được trống")
	}
	validStatuses := map[string]bool{"pending": true, "in_progress": true, "done": true}
	if !validStatuses[todo.Status] {
		return errors.New("status không hợp lệ(chỉ chấp nhận : pending , in_progress , done)")
	}

	return s.repo.Update(todo)
}

func (s *TodoService) Delete(id uint) error {
	return s.repo.Delete(id)
}
