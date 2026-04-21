package repositories

import (
	"github.com/DoDtatt/todo-app/internal/models"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *TodoRepository) GetAll(userID int, status string, search string, sort string, order string) ([]models.Todo, error) {
	var todos []models.Todo

	query := r.db.Model(&models.Todo{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	query = query.Order(sort + " " + order)

	err := query.Find(&todos).Error
	return todos, err
}

func (r *TodoRepository) GetbyID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	return &todo, err
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.db.Updates(todo).Error
}

func (r *TodoRepository) Delete(id uint) error {
	var todo models.Todo
	return r.db.Delete(&todo, id).Error
}
