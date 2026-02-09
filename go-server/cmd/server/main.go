package main

import (
	"fmt"
	"log"
	"os"
	"server/internal/auth"
	"server/internal/config"
	"server/internal/resource"
	"server/internal/role"
	"server/internal/router"
	"server/internal/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	// 1. Try current directory first (specific config) -> won't be overwritten by subsequent loads
	if err := godotenv.Load(); err != nil {
		// 2. If not found or just to fill gaps, try parent directory
		godotenv.Load("../.env")
	} else {
		// Even if found, try parent to fill any missing global vars
		godotenv.Load("../.env")
	}

	// Initialize database
	config.InitDB()

	// Initialize Gin
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Initialize repositories
	authRepo := &auth.Repository{}
	userRepo := &user.Repository{}
	roleRepo := &role.Repository{}
	resourceRepo := &resource.Repository{}

	// Initialize services
	authService := &auth.Service{Repo: authRepo}
	userService := &user.Service{Repo: userRepo}
	roleService := &role.Service{Repo: roleRepo}
	resourceService := &resource.Service{Repo: resourceRepo}

	// Initialize handlers
	authHandler := &auth.Handler{Service: authService}
	userHandler := &user.Handler{Service: userService}
	roleHandler := &role.Handler{Service: roleService}
	resourceHandler := &resource.Handler{Service: resourceService}

	// Setup routes
	router.SetupRoutes(r, authHandler, userHandler, roleHandler, resourceHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
