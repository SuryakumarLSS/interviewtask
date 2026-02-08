package handlers

import (
	"net/http"
	"server/services"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	ResourceService *services.ResourceService
}

func (h *ResourceHandler) GetAll(c *gin.Context) {
	resource := c.Param("resource")
	// Check Access
	// Original code used middleware.RBACMiddleware
	// We should continue to use middleware or check here.
	// If the user wants middleware to do it, we keep middleware.
	// If we move logic here:
	// user, _ := c.Get("user"); claims := user.(*utils.Claims); allowed, _ := h.ResourceService.CheckAccess(claims.RoleID, resource, "read")
	// But original code applied middleware to the *group*.
	// I will keep the middleware approach in router Setup, so this handler is just business.
	// Wait, original code applied middleware INSIDE the route handler manually in data.go:
	// middleware.RBACMiddleware(resource, "read")(c)
	// I should replicate that pattern or clean it up.
	// I'll use the service check here to be cleaner "Layered".

	user, _ := c.Get("user")
	claims := user.(*utils.Claims)
	allowed, err := h.ResourceService.CheckAccess(claims.RoleID, resource, "read")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
		return
	}
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"message": "Access Denied"})
		return
	}

	data, err := h.ResourceService.GetAll(resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ResourceHandler) Create(c *gin.Context) {
	resource := c.Param("resource")

	user, _ := c.Get("user")
	claims := user.(*utils.Claims)
	allowed, _ := h.ResourceService.CheckAccess(claims.RoleID, resource, "create")
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"message": "Access Denied"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.ResourceService.Create(resource, data)
	if err != nil {
		if strings.Contains(err.Error(), "Missing required fields") {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	data["id"] = id
	c.JSON(http.StatusOK, data)
}

func (h *ResourceHandler) Update(c *gin.Context) {
	resource := c.Param("resource")
	id := c.Param("id")

	user, _ := c.Get("user")
	claims := user.(*utils.Claims)
	allowed, _ := h.ResourceService.CheckAccess(claims.RoleID, resource, "update")
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"message": "Access Denied"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ResourceService.Update(resource, id, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func (h *ResourceHandler) Delete(c *gin.Context) {
	resource := c.Param("resource")
	id := c.Param("id")

	user, _ := c.Get("user")
	claims := user.(*utils.Claims)
	allowed, _ := h.ResourceService.CheckAccess(claims.RoleID, resource, "delete")
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"message": "Access Denied"})
		return
	}

	err := h.ResourceService.Delete(resource, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// Special Gets

func (h *ResourceHandler) GetEmployees(c *gin.Context) {
	// Helper route from auth.go
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No token"})
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	data, err := h.ResourceService.GetEmployeeNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ResourceHandler) GetUsers(c *gin.Context) {
	// Helper route from auth.go
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No token"})
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	data, err := h.ResourceService.GetUserList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ResourceHandler) GetMyPermissions(c *gin.Context) {
	// Helper route from auth.go
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No token"})
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	perms, err := h.ResourceService.GetPermissions(claims.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}
