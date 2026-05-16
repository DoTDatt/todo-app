package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/repositories"
	"github.com/DoDtatt/todo-app/internal/search"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm"
)

type TodoService struct {
	repo  *repositories.TodoRepository
	meili *search.Meili
}

func NewtodoService(repo *repositories.TodoRepository, meili *search.Meili) *TodoService {
	return &TodoService{repo: repo, meili: meili}
}

func (s *TodoService) CreateTodo(todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("Title không được trống")
	}
	validStatuses := map[string]bool{"pending": true, "in_progress": true, "done": true}
	if !validStatuses[todo.Status] {
		return errors.New("status không hợp lệ(chỉ chấp nhận : pending , in_progress , done)")
	}

	tx := s.repo.DB().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := s.repo.Create(tx, todo); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	if err := s.meili.UpsertDocuments("todos", todo); err != nil {
		fmt.Printf("[Meili] UpsertDocument failed: %v\n", err)
	} else {
		fmt.Printf("[Meili OK] synced todo id=%d\n", todo.ID)
	}
	return nil
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

	tx := s.repo.DB().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := s.repo.Update(tx, todo); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	if err := s.meili.UpsertDocuments("todos", todo); err != nil {

		fmt.Printf("[Meili] UpsertDocument failed: %v\n", err)
	} else {
		fmt.Printf("[Meili OK] synced todo id=%d\n", todo.ID)
	}
	return nil

}

func (s *TodoService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if err := s.meili.DeleteDocument("todos", id); err != nil {
		fmt.Printf("[Meili] DeleteDocument failed: %v\n", err)
	}
	return nil
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

func (s *TodoService) SearchMeili(query string) ([]models.Todo, error) {
	res, err := s.meili.Client.Index("todos").Search(query, &meilisearch.SearchRequest{
		Limit: 20,
	})
	if err != nil {
		return nil, err
	}

	var todos []models.Todo
	for _, hit := range res.Hits {
		data, _ := json.Marshal(hit)

		var raw map[string]interface{}
		json.Unmarshal(data, &raw)

		if t, ok := raw["created_at"]; ok {
			raw["created_at"] = fmt.Sprintf("%v", t)
		}
		if t, ok := raw["updated_at"]; ok {
			raw["updated_at"] = fmt.Sprintf("%v", t)
		}

		fixed, _ := json.Marshal(raw)
		var todo models.Todo
		json.Unmarshal(fixed, &todo)
		todos = append(todos, todo)
	}
	return todos, nil
}
