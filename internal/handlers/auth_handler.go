package handlers

import (
	"net/http"

	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	serv *services.AuthService
}

func NewAuthHandler(serv *services.AuthService) *AuthHandler {
	return &AuthHandler{serv: serv}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.serv.Register(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Đăng ký thành công"})

}
