package handlers

import (
	"fmt"
	"net/http"
	"server/services"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminService *services.AdminService
}

func (h *AdminHandler) GetRoles(c *gin.Context) {
	roles, err := h.AdminService.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	// Original code bind to role struct but only used name
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.AdminService.CreateRole(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "name": req.Name})
}

func (h *AdminHandler) GetPermissions(c *gin.Context) {
	perms, err := h.AdminService.GetPermissions(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}

func (h *AdminHandler) AddOrUpdatePermission(c *gin.Context) {
	var req struct {
		RoleID   int    `json:"role_id"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, id, err := h.AdminService.AddOrUpdatePermission(req.RoleID, req.Resource, req.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if status == "updated" {
		c.JSON(http.StatusOK, gin.H{"message": "Permission updated to full access"})
	} else {
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "Permission created"})
	}
}

func (h *AdminHandler) DeletePermission(c *gin.Context) {
	roleID := c.Query("role_id")
	resource := c.Query("resource")
	action := c.Query("action")

	err := h.AdminService.DeletePermission(roleID, resource, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Permission removed"})
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	users, err := h.AdminService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RoleID   int    `json:"role_id"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, link, id, err := h.AdminService.CreateUserOrInvite(struct {
		Username string
		Password string
		RoleID   int
		Email    string
	}{
		Username: req.Username,
		Password: req.Password,
		RoleID:   req.RoleID,
		Email:    req.Email,
	})
	if err != nil {
		// status maps to error sometimes in service? No service returns err separately.
		// If custom error message needed:
		if err.Error() == "User already registered" {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if status == "invite" {
		c.JSON(http.StatusOK, gin.H{"message": "Invitation sent (Mock)", "link": link})
	} else {
		c.JSON(http.StatusOK, gin.H{"id": id, "username": req.Username, "role_id": req.RoleID})
	}
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	// Superadmin check
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)
	if claims.Username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only Superadmin can perform this action"})
		return
	}

	var req struct {
		RoleID int `json:"role_id"`
	}
	c.ShouldBindJSON(&req)

	// Since we need to parse ID from param string, let's do simple validation or conversion if needed.
	// But DB exec usually handles string params fine for integer columns in current driver?
	// The original code passed c.Param("id") directly to DB.Exec.
	// But `AdminService.UpdateUserRole` takes int.
	// Let's rely on standard binding or parse.
	// Actually for simplicity, let `UpdateUserRole` take interface{} or we parse.
	// I'll parse it here.
	var id int
	fmt.Sscan(c.Param("id"), &id)

	err := h.AdminService.UpdateUserRole(id, req.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	// Superadmin check
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)
	if claims.Username != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only Superadmin can perform this action"})
		return
	}

	err := h.AdminService.DeleteUser(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
