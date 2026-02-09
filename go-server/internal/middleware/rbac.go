package middleware

import (
	"server/internal/config"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatus(401)
			return
		}

		claims := user.(*utils.Claims)

		// Admin role (assuming ID 1 is Admin, or check name)
		// Optimized: If Admin, maybe always allow? But let's follow DB rules as per node app.
		if claims.RoleID == 1 {
			c.Next()
			return
		}

		// Check permission using new RBAC schema
		query := `
			SELECT CASE 
				WHEN $3 = 'read' THEN rrp.can_view
				WHEN $3 = 'create' THEN rrp.can_create
				WHEN $3 = 'update' THEN rrp.can_update
				WHEN $3 = 'delete' THEN rrp.can_delete
				ELSE false
			END as has_permission
			FROM role_resource_permissions rrp
			JOIN resources res ON rrp.resource_id = res.id
			WHERE rrp.role_id = $1 AND res.name = $2
		`

		var hasPermission bool
		err := config.DB.QueryRow(query, claims.RoleID, resource, action).Scan(&hasPermission)

		if err != nil || !hasPermission {
			c.JSON(403, gin.H{"message": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
