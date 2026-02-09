package user

import (
	"fmt"
	"net/http"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func (h *Handler) GetAll(c *gin.Context) {
	users, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, link, id, err := h.Service.CreateOrInvite(req)
	if err != nil {
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
		c.JSON(http.StatusOK, gin.H{"id": id, "username": req.Email, "role_id": req.RoleID})
	}
}

func (h *Handler) UpdateRole(c *gin.Context) {
	// Superadmin check using is_admin flag
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)

	isAdmin, err := h.Service.Repo.IsAdmin(claims.ID)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only Superadmin can perform this action"})
		return
	}

	var req UpdateUserRoleRequest
	c.ShouldBindJSON(&req)

	var id int
	fmt.Sscan(c.Param("id"), &id)

	err = h.Service.UpdateRole(id, req.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	// Superadmin check using is_admin flag
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)

	isAdmin, err := h.Service.Repo.IsAdmin(claims.ID)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only Superadmin can perform this action"})
		return
	}

	err = h.Service.Delete(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
