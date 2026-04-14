package handlers

import (
	"fmt"
	"net/http"

	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/repositories"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	repo *repositories.TodoRepository
}

func NewTodoHandler(repo *repositories.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	todo, err := h.repo.GetbyID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todo models.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lưu dữ liệu"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}
