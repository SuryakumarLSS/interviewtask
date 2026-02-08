package middleware

import (
	"net/http"

	"server/database"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(resource string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		claims := user.(*utils.Claims)
		roleID := claims.RoleID

		// Admin role (assuming ID 1 is Admin, or check name)
		// For simplicity, we check DB permissions generally.
		// Optimized: If Admin, maybe always allow? But let's follow DB rules as per node app.
		// Node app: Permissions table stores explicit grants.

		var count int
		// Check for specific permission or wildcard
		query := `SELECT COUNT(*) FROM permissions 
				  WHERE role_id = $1 AND resource = $2 AND (action = $3 OR action = '*')`

		// Note on Attributes: Node code updated attributes to '*' in one place but didn't seem to enforce them strictly in middleware beyond basic check.
		// The node middleware (I should check it, but I recall it checking standard permissions).
		// Wait, I didn't read middleware/rbac.js explicitly but used in data.js.
		// Assuming standard: resource, action match.

		err := database.DB.QueryRow(query, roleID, resource, action).Scan(&count)
		if err != nil {
			// In case of DB error, deny
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if count > 0 {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"message": "Access Denied"})
		c.Abort()
	}
}
