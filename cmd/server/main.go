package main

import (
	"github.com/DoDtatt/todo-app/internal/config"
	"github.com/DoDtatt/todo-app/internal/handlers"
	"github.com/DoDtatt/todo-app/internal/middleware"
	"github.com/DoDtatt/todo-app/internal/models"
	"github.com/DoDtatt/todo-app/internal/repositories"
	"github.com/DoDtatt/todo-app/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()
	config.DB.AutoMigrate(&models.Todo{})

	todorepo := repositories.NewTodoRepository(config.DB)
	todoserv := services.NewtodoService(todorepo)
	todoHandler := handlers.NewTodoHandler(todoserv)

	authRepo := repositories.NewAuthRepository(config.DB)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	r.Use(middleware.RecoverMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RateLimit())
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	protected := r.Group("/api")

	protected.GET("/todos", middleware.AuthMiddleware(), todoHandler.GetAllTodos)
	protected.POST("/todos", middleware.AuthMiddleware(), todoHandler.CreateTodo)
	protected.GET("/todos/:id", middleware.AuthMiddleware(), todoHandler.GetTodoByID)
	protected.PUT("/todos/:id", middleware.AuthMiddleware(), todoHandler.UpdateTodo)
	protected.DELETE("/todos/:id", middleware.AuthMiddleware(), todoHandler.DeleteTodo)

	r.Run(":8080")
}
