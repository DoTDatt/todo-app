package services

import (
	"errors"
	"fmt"

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

func (s *TodoService) GetbyID(id int) (*models.Todo, error) {
	todo, err := s.repo.GetbyID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Todo không tồn tại")
	}
	return todo, nil
}

func (s *TodoService) GetAll(userID int, p models.TodoQuery) ([]models.Todo, error) {

	allowedSort := map[string]bool{"id": true, "title": true, "created_at": true}

	if p.Sort == "" {
		p.Sort = "id"
	}

	if !allowedSort[p.Sort] {
		p.Sort = "id"
	}

	if p.Order != "asc" && p.Order != "desc" {
		p.Order = "desc"
	}

	if p.Status != "" {
		validStatuses := map[string]bool{"pending": true, "in_progress": true, "done": true}
		if !validStatuses[p.Status] {
			return nil, errors.New("status chỉ chấp nhận pending, in_progress hoặc done")
		}
	}

	scopes := []repositories.Scope{
		s.Status(p.Status),
		s.Search(p.Search),
		s.Sort(p.Sort, p.Order),
	}

	return s.repo.GetAll(userID, scopes...)
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

func (s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TodoService) Status(status string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status != "" {
			return db.Where("status = ?", status)
		}
		return db
	}
}

func (s *TodoService) Search(search string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search != "" {
			return db.Where("title LIKE ?", "%"+search+"%")
		}
		return db
	}
}

func (s *TodoService) Sort(sort, order string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		allowedColumns := map[string]bool{"id": true, "title": true, "created_at": true}
		if !allowedColumns[sort] {
			sort = "id"
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		return db.Order(fmt.Sprintf("%s %s", sort, order))
	}
}
