package repositories

import (
	"github.com/DoDtatt/todo-app/internal/models"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

type Scope func(*gorm.DB) *gorm.DB

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *TodoRepository) GetAll(userID int, scopes ...Scope) ([]models.Todo, error) {
	var todos []models.Todo
	db := r.db.Model(&models.Todo{}).Where("user_id = ?", userID)
	for _, s := range scopes {
		db = db.Scopes(s)
	}
	err := db.Find(&todos).Error

	return todos, err
}

func (r *TodoRepository) GetbyID(id int) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	return &todo, err
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.db.Updates(todo).Error
}

func (r *TodoRepository) Delete(id int) error {
	var todo models.Todo
	return r.db.Delete(&todo, id).Error
}
