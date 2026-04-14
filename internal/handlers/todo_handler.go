package handlers

import (
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
