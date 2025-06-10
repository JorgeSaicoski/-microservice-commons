package main

import (
	"time"

	"github.com/JorgeSaicoski/microservice-commons/config"
	"github.com/JorgeSaicoski/microservice-commons/database"
	"github.com/JorgeSaicoski/microservice-commons/middleware"
	"github.com/JorgeSaicoski/microservice-commons/responses"
	"github.com/JorgeSaicoski/microservice-commons/server"
	"github.com/JorgeSaicoski/microservice-commons/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Project model with full features
type Project struct {
	types.BaseModel
	Name        string         `json:"name" gorm:"not null;index"`
	Description string         `json:"description"`
	Status      types.Status   `json:"status" gorm:"default:'active'"`
	Priority    types.Priority `json:"priority" gorm:"default:'medium'"`
	OwnerID     uint           `json:"owner_id" gorm:"not null;index"`
	Tags        types.Tags     `json:"tags" gorm:"type:text[]"`
	Metadata    types.Metadata `json:"metadata" gorm:"type:jsonb"`
}

// User model for demonstration
type User struct {
	types.BaseModel
	Username string           `json:"username" gorm:"uniqueIndex;not null"`
	Email    string           `json:"email" gorm:"uniqueIndex;not null"`
	Name     string           `json:"name"`
	Status   types.UserStatus `json:"status" gorm:"default:'active'"`
	Role     types.Role       `json:"role" gorm:"default:'user'"`
}

var db *gorm.DB

func main() {
	// Create server with advanced configuration
	server := server.NewServer(server.ServerOptions{
		ServiceName:    "advanced-project-service",
		ServiceVersion: "2.0.0",
		SetupRoutes:    setupRoutes,
		// Custom middleware
		CustomMiddleware: []gin.HandlerFunc{
			middleware.DefaultRequestIDMiddleware(),
			middleware.RequestLogger(middleware.LogLevelInfo),
		},
	})

	server.Start()
}

func setupRoutes(router *gin.Engine, cfg *config.Config) {
	// Connect to database
	var err error
	db, err = database.ConnectWithConfig(cfg.DatabaseConfig)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Migrate with custom options
	migrator := database.NewMigrator(db, database.MigrationOptions{
		DropTables:    false,
		CreateIndexes: true,
		Verbose:       cfg.IsDevelopment(),
	})

	if err := migrator.AddModels(&User{}, &Project{}).Migrate(); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Setup health checks with database monitoring
	setupHealthChecks(router, cfg)

	// Setup authentication (mock for example)
	authMiddleware := setupAuthentication()

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/users/register", registerUser)
		public.POST("/users/login", loginUser)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware)
	{
		// Users
		users := protected.Group("/users")
		{
			users.GET("", getUsers)
			users.GET("/:id", getUser)
			users.PUT("/:id", updateUser)
		}

		// Projects
		projects := protected.Group("/projects")
		{
			projects.GET("", getProjects) // With pagination
			projects.POST("", createProject)
			projects.GET("/:id", getProject)
			projects.PUT("/:id", updateProject)
			projects.DELETE("/:id", deleteProject)
			projects.GET("/search", searchProjects) // With filtering
		}
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(authMiddleware)
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/stats", getStatistics)
		admin.DELETE("/users/:id", deleteUser)
	}
}

func setupHealthChecks(router *gin.Engine, cfg *config.Config) {
	// Custom health check with database monitoring
	healthConfig := middleware.DefaultHealthConfig("advanced-project-service", "2.0.0")

	// Add database health checker
	healthConfig.AddHealthChecker("database", middleware.DatabaseHealthChecker(func() error {
		return database.QuickHealthCheck(db)
	}))

	// Add memory health checker (max 512MB)
	healthConfig.AddHealthChecker("memory", middleware.MemoryHealthChecker(512))

	router.Use(middleware.HealthMiddleware(healthConfig))
}

func setupAuthentication() gin.HandlerFunc {
	// Mock authentication - in real world, use JWT/Keycloak
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			responses.Unauthorized(c, "Authorization token required")
			c.Abort()
			return
		}

		// Mock user data
		c.Set("user_id", "1")
		c.Set("username", "testuser")
		c.Set("roles", []string{"user"})

		// For admin endpoints demo
		if token == "Bearer admin-token" {
			c.Set("roles", []string{"admin"})
		}

		c.Next()
	}
}

// User handlers
func registerUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		responses.BadRequestWithDetails(c, "Invalid user data", err.Error())
		return
	}

	// Validate user data
	if !isValidEmail(user.Email) {
		responses.BadRequest(c, "Invalid email format")
		return
	}

	if err := db.Create(&user).Error; err != nil {
		responses.Conflict(c, "User already exists")
		return
	}

	responses.Created(c, "User registered successfully", user)
}

func loginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		responses.BadRequest(c, "Invalid credentials")
		return
	}

	// Mock login response
	response := map[string]interface{}{
		"token":      "mock-jwt-token",
		"expires_in": 3600,
		"user": map[string]interface{}{
			"id":       1,
			"username": "testuser",
			"email":    credentials.Email,
		},
	}

	responses.Success(c, "Login successful", response)
}

func getUsers(c *gin.Context) {
	// Get pagination parameters
	params := responses.GetPaginationParams(c)

	var users []User
	var total int64

	// Count total records
	db.Model(&User{}).Count(&total)

	// Get paginated results
	if err := db.Offset(params.Offset).Limit(params.Limit).Find(&users).Error; err != nil {
		responses.InternalError(c, "Failed to fetch users")
		return
	}

	responses.Paginated(c, users, total, params.Page, params.PageSize)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.NotFound(c, "User not found")
			return
		}
		responses.InternalError(c, "Failed to fetch user")
		return
	}

	responses.Success(c, "User retrieved successfully", user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.NotFound(c, "User not found")
			return
		}
		responses.InternalError(c, "Failed to fetch user")
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		responses.BadRequest(c, "Invalid user data")
		return
	}

	if err := db.Save(&user).Error; err != nil {
		responses.InternalError(c, "Failed to update user")
		return
	}

	responses.Success(c, "User updated successfully", user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	result := db.Delete(&User{}, id)
	if result.Error != nil {
		responses.InternalError(c, "Failed to delete user")
		return
	}

	if result.RowsAffected == 0 {
		responses.NotFound(c, "User not found")
		return
	}

	responses.Success(c, "User deleted successfully", nil)
}

// Project handlers with advanced features
func getProjects(c *gin.Context) {
	params := responses.GetPaginationParams(c)

	// Advanced filtering
	query := db.Model(&Project{})

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by priority
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Search by name
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	var projects []Project
	var total int64

	query.Count(&total)

	if err := query.Offset(params.Offset).Limit(params.Limit).Find(&projects).Error; err != nil {
		responses.InternalError(c, "Failed to fetch projects")
		return
	}

	responses.Paginated(c, projects, total, params.Page, params.PageSize)
}

func createProject(c *gin.Context) {
	var project Project

	if err := c.ShouldBindJSON(&project); err != nil {
		responses.BadRequest(c, "Invalid project data")
		return
	}

	// Set owner from auth context
	userID, _ := c.Get("user_id")
	project.OwnerID = 1 // Mock user ID

	// Set default metadata
	project.Metadata = types.Metadata{
		"created_by_service": "advanced-project-service",
		"version":            "2.0.0",
	}

	if err := db.Create(&project).Error; err != nil {
		responses.InternalError(c, "Failed to create project")
		return
	}

	responses.Created(c, "Project created successfully", project)
}

func getProject(c *gin.Context) {
	id := c.Param("id")
	var project Project

	if err := db.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.NotFound(c, "Project not found")
			return
		}
		responses.InternalError(c, "Failed to fetch project")
		return
	}

	responses.Success(c, "Project retrieved successfully", project)
}

func updateProject(c *gin.Context) {
	id := c.Param("id")
	var project Project

	if err := db.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.NotFound(c, "Project not found")
			return
		}
		responses.InternalError(c, "Failed to fetch project")
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		responses.BadRequest(c, "Invalid project data")
		return
	}

	if err := db.Save(&project).Error; err != nil {
		responses.InternalError(c, "Failed to update project")
		return
	}

	responses.Success(c, "Project updated successfully", project)
}

func deleteProject(c *gin.Context) {
	id := c.Param("id")

	result := db.Delete(&Project{}, id)
	if result.Error != nil {
		responses.InternalError(c, "Failed to delete project")
		return
	}

	if result.RowsAffected == 0 {
		responses.NotFound(c, "Project not found")
		return
	}

	responses.Success(c, "Project deleted successfully", nil)
}

func searchProjects(c *gin.Context) {
	searchTerm := c.Query("q")
	if searchTerm == "" {
		responses.BadRequest(c, "Search query is required")
		return
	}

	var projects []Project

	if err := db.Where("name ILIKE ? OR description ILIKE ?",
		"%"+searchTerm+"%", "%"+searchTerm+"%").Find(&projects).Error; err != nil {
		responses.InternalError(c, "Search failed")
		return
	}

	response := map[string]interface{}{
		"query":   searchTerm,
		"results": projects,
		"count":   len(projects),
	}

	responses.Success(c, "Search completed", response)
}

func getStatistics(c *gin.Context) {
	var userCount, projectCount int64

	db.Model(&User{}).Count(&userCount)
	db.Model(&Project{}).Count(&projectCount)

	stats := map[string]interface{}{
		"users":    userCount,
		"projects": projectCount,
		"uptime":   time.Since(time.Now()).String(),
	}

	responses.Success(c, "Statistics retrieved", stats)
}

// Utility functions
func isValidEmail(email string) bool {
	// Basic email validation - in real app use proper validation
	return len(email) > 5 && email[len(email)-4:] == ".com"
}
