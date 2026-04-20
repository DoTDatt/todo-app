package repositories

import (
	"github.com/DoDtatt/todo-app/internal/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Raw("SELECT * FROM users WHERE email = ? LIMIT 1", email).Scan(&user).Error
	return &user, err
}

func (r *AuthRepository) GetRoleIDByName(name string) (int, error) {
	var id int
	err := r.DB.Raw("SELECT id FROM roles WHERE role_name = ?", name).Scan(&id).Error
	return id, err
}

func (r *AuthRepository) CreateUser(email, password string, roleID int) error {
	return r.DB.Exec("INSERT INTO users (email, password, role_id) VALUES (?, ?, ?)", email, password, roleID).Error
}
