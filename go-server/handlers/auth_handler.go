package handlers

import (
	"net/http"
	"server/services"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *services.AuthService
	UserRepo    map[string]interface{} // Hack: Actually we need full wiring
	// But let's use methods.
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.AuthService.Login(loginData.Username, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "role_id": user.RoleID},
	})
}

func (h *AuthHandler) SetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AuthService.SetPassword(req.Token, req.Password); err != nil {
		// Matching error message from original code roughly
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password set successfully"})
}

func (h *AuthHandler) DeclineInvitation(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AuthService.DeclineInvitation(req.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation declined"})
}

// GetPermissions handles GET /auth/permissions
func (h *AuthHandler) GetPermissions(c *gin.Context, roleRepo interface{}) {
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

	// We need RoleService/Repo here.
	// Since Handler struct definition isn't final in this snippet, I'll assume we pass it or attach it.
	// Let's assume h struct has a service for this.
}
