package handlers

import (
	"fmt"
	"net/http"

	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/services"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	serv *services.TodoService
}

func NewTodoHandler(serv *services.TodoService) *TodoHandler {
	return &TodoHandler{serv: serv}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todo models.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := c.Get("user_id")

	todo.UserID = uid.(int)
	if err := h.serv.CreateTodo(&todo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	todo, err := h.serv.GetbyID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi lấy todo"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.serv.GetAll(c.GetInt("user_id"), c.Query("status"), c.Query("search"), c.Query("sort"), c.Query("order"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi lấy danh sách todo"})
		return
	}
	if len(todos) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Không có todo nào"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	idParam := c.Param("id")
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID không hợp lệ"})
		return
	}
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo.ID = id
	if err := h.serv.Update(&todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.serv.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo đã được xóa"})
}
