package main

import (
	"github.com/DoDtatt/todo-app/internal/config"
	"github.com/DoDtatt/todo-app/internal/handlers"
	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/repositories"
	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()
	config.DB.AutoMigrate(&models.Todo{})

	r := gin.Default()

	todorepo := repositories.NewTodoRepository(config.DB)
	todoHandler := handlers.NewTodoHandler(todorepo)

	r.POST("/todos", todoHandler.CreateTodo)
	r.GET("/todos", todoHandler.GetAllTodos)
	r.GET("/todos/:id", todoHandler.GetTodoByID)
	r.Run(":8080")
}
