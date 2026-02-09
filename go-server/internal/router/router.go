package router

import (
	"server/internal/auth"
	"server/internal/middleware"
	"server/internal/resource"
	"server/internal/role"
	"server/internal/user"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *auth.Handler,
	userHandler *user.Handler,
	roleHandler *role.Handler,
	resourceHandler *resource.Handler,
) {
	api := r.Group("/api")

	// Auth routes (public + authenticated)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/set-password", authHandler.SetPassword)
		authGroup.POST("/decline-invitation", authHandler.DeclineInvitation)

		// Authenticated auth routes
		authenticatedAuth := authGroup.Group("/")
		authenticatedAuth.Use(middleware.AuthMiddleware())
		{
			authenticatedAuth.GET("/permissions", authHandler.GetMyPermissions)
			authenticatedAuth.GET("/field-permissions", authHandler.GetMyFieldPermissions)
			authenticatedAuth.GET("/employees", authHandler.GetEmployees)
			authenticatedAuth.GET("/users", authHandler.GetUsers)
		}
	}

	// Admin routes (Admin role only)
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(func(c *gin.Context) {
		// Admin role check (RoleID == 1)
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
		// Role management
		adminGroup.GET("/roles", roleHandler.GetAll)
		adminGroup.POST("/roles", roleHandler.Create)
		adminGroup.DELETE("/roles/:id", roleHandler.Delete)
		adminGroup.GET("/permissions/:role_id", roleHandler.GetPermissions)
		adminGroup.POST("/permissions", roleHandler.AddOrUpdatePermission)
		adminGroup.DELETE("/permissions", roleHandler.DeletePermission)

		// Field-level permissions
		adminGroup.GET("/field-permissions/:role_id", roleHandler.GetFieldPermissions)
		adminGroup.POST("/field-permissions", roleHandler.UpdateFieldPermission)

		// User management
		adminGroup.GET("/users", userHandler.GetAll)
		adminGroup.POST("/users", userHandler.Create)
		adminGroup.PUT("/users/:id/role", userHandler.UpdateRole)
		adminGroup.DELETE("/users/:id", userHandler.Delete)
	}

	// Data routes (Authenticated with resource validation)
	dataGroup := api.Group("/data")
	dataGroup.Use(middleware.AuthMiddleware())
	{
		// Resource validation middleware
		dataGroup.Use(func(c *gin.Context) {
			resource := c.Param("resource")
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
}
