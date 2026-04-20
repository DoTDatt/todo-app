package services

import (
	"errors"

	"github.com/DoDtatt/todo-app/internal/repositories"
	"github.com/DoDtatt/todo-app/internal/utils"
)

type AuthService struct {
	Repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(email, password string) error {

	user, _ := s.Repo.FindByEmail(email)
	if user.ID != 0 {
		return errors.New("email đã tồn tại")
	}

	roleID, err := s.Repo.GetRoleIDByName("user")
	if err != nil || roleID == 0 {
		return errors.New("không tìm thấy role mặc định")
	}

	hashed, _ := utils.HashPassword(password)

	return s.Repo.CreateUser(email, hashed, roleID)
}

func (s *AuthService) Login(email, password string) (string, error) {

	user, err := s.Repo.FindByEmail(email)
	if err != nil || user.ID == 0 {
		return "", errors.New("email không tồn tại")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("mật khẩu không chính xác")
	}

	token, err := utils.GenerateJWT(int(user.ID), user.RoleID)
	if err != nil {
		return "", errors.New("lỗi khi tạo token")
	}

	return token, nil
}
