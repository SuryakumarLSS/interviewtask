package main

import (
	"log"
	"os"

	"server/database"
	"server/handlers"
	"server/middleware"
	"server/repositories"
	"server/services"
	"server/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env from current or parent directory
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found")
		}
	}

	database.InitDB()

	// Dependency Injection Setup

	// Repositories
	userRepo := &repositories.UserRepository{}
	roleRepo := &repositories.RoleRepository{}
	resourceRepo := &repositories.ResourceRepository{}

	// Services
	authService := &services.AuthService{UserRepo: userRepo}
	adminService := &services.AdminService{UserRepo: userRepo, RoleRepo: roleRepo}
	resourceService := &services.ResourceService{
		ResourceRepo: resourceRepo,
		UserRepo:     userRepo,
		RoleRepo:     roleRepo,
	}

	// Handlers
	authHandler := &handlers.AuthHandler{AuthService: authService}
	adminHandler := &handlers.AdminHandler{AdminService: adminService}
	resourceHandler := &handlers.ResourceHandler{ResourceService: resourceService}

	r := gin.Default()

	// CORS Setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	api := r.Group("/api")

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/set-password", authHandler.SetPassword)
		authGroup.POST("/decline-invitation", authHandler.DeclineInvitation)
		// Old auth.go mixed handlers, now mapped to ResourceHandler/AuthHandler as appropriate
		authGroup.GET("/permissions", resourceHandler.GetMyPermissions)
		authGroup.GET("/employees", resourceHandler.GetEmployees)
		authGroup.GET("/users", resourceHandler.GetUsers) // Authenticated list of users
	}

	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	// Additional Admin Check Middleware - simplified to check inside handler or use custom
	adminGroup.Use(func(c *gin.Context) {
		// We could use the old check here or inside handlers.
		// Handlers have checks but middleware is safer.
		// For now, let's keep the middleware logic inline or re-implement if needed.
		// admin_routes.go had: r.Use(...) check ID==1.
		// Let's add it back inline here for safety.
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatus(401)
			return
		}

		claims := user.(*utils.Claims)
		if claims.RoleID != 1 {
			c.JSON(403, gin.H{"message": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	})

	{
		adminGroup.GET("/roles", adminHandler.GetRoles)
		adminGroup.POST("/roles", adminHandler.CreateRole)
		adminGroup.GET("/permissions/:role_id", adminHandler.GetPermissions)
		adminGroup.POST("/permissions", adminHandler.AddOrUpdatePermission)
		adminGroup.DELETE("/permissions", adminHandler.DeletePermission)
		adminGroup.GET("/users", adminHandler.GetUsers)
		adminGroup.POST("/users", adminHandler.CreateUser)
		adminGroup.PUT("/users/:id/role", adminHandler.UpdateUserRole)
		adminGroup.DELETE("/users/:id", adminHandler.DeleteUser)
	}

	dataGroup := api.Group("/data")
	dataGroup.Use(middleware.AuthMiddleware())
	{
		// check resource validity middleware inline
		dataGroup.Use(func(c *gin.Context) {
			// Logic from data.go
			resource := c.Param("resource")
			// We need AllowedResources map. It was private in routes package.
			// Let's just check against the hardcoded list or move map to utils/config.
			// Or check if in our allowed list.
			allowed := map[string]bool{"employees": true, "projects": true, "orders": true}
			if !allowed[resource] {
				c.JSON(404, gin.H{"message": "Resource not found"})
				c.Abort()
				return
			}
			c.Next()
		})

		dataGroup.GET("/:resource", resourceHandler.GetAll)
		dataGroup.POST("/:resource", resourceHandler.Create)
		dataGroup.PUT("/:resource/:id", resourceHandler.Update)
		dataGroup.DELETE("/:resource/:id", resourceHandler.Delete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
