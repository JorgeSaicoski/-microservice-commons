package main

import (
	"github.com/JorgeSaicoski/microservice-commons/config"
	"github.com/JorgeSaicoski/microservice-commons/database"
	"github.com/JorgeSaicoski/microservice-commons/responses"
	"github.com/JorgeSaicoski/microservice-commons/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Simple task model for demonstration
type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Completed   bool   `json:"completed" gorm:"default:false"`
}

var db *gorm.DB

func main() {
	// Create server with microservice-commons
	server := server.NewServer(server.ServerOptions{
		ServiceName:    "basic-task-service",
		ServiceVersion: "1.0.0",
		SetupRoutes:    setupRoutes,
	})

	// Start the server (includes graceful shutdown)
	server.Start()
}

func setupRoutes(router *gin.Engine, cfg *config.Config) {
	// Connect to database using microservice-commons
	var err error
	db, err = database.ConnectWithConfig(cfg.DatabaseConfig)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Auto-migrate the task model
	if err := database.QuickMigrate(db, &Task{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Setup API routes
	api := router.Group("/api/v1")
	{
		api.GET("/tasks", getTasks)
		api.POST("/tasks", createTask)
		api.GET("/tasks/:id", getTask)
		api.PUT("/tasks/:id", updateTask)
		api.DELETE("/tasks/:id", deleteTask)
	}
}

func getTasks(c *gin.Context) {
	var tasks []Task

	if err := db.Find(&tasks).Error; err != nil {
		responses.InternalError(c, "Failed to fetch tasks")
		return
	}

	responses.Success(c, "Tasks retrieved successfully", tasks)
}

func createTask(c *gin.Context) {
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		responses.BadRequest(c, "Invalid task data")
		return
	}

	if err := db.Create(&task).Error; err != nil {
		responses.InternalError(c, "Failed to create task")
		return
	}

	responses.Created(c, "Task created successfully", task)
}

func getTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.NotFound(c, "Task not found")
			return
		}
		responses.InternalError(c, "Failed to fetch task")
		return
	}

	responses.Success(c, "Task retrieved successfully", task)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.NotFound(c, "Task not found")
			return
		}
		responses.InternalError(c, "Failed to fetch task")
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		responses.BadRequest(c, "Invalid task data")
		return
	}

	if err := db.Save(&task).Error; err != nil {
		responses.InternalError(c, "Failed to update task")
		return
	}

	responses.Success(c, "Task updated successfully", task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	result := db.Delete(&Task{}, id)
	if result.Error != nil {
		responses.InternalError(c, "Failed to delete task")
		return
	}

	if result.RowsAffected == 0 {
		responses.NotFound(c, "Task not found")
		return
	}

	responses.Success(c, "Task deleted successfully", nil)
}
