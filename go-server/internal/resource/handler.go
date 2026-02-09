package resource

import (
	"net/http"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func (h *Handler) GetAll(c *gin.Context) {
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)

	resource := c.Param("resource")
	data, err := h.Service.GetAll(resource, claims.RoleID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) Create(c *gin.Context) {
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)

	resource := c.Param("resource")
	var data map[string]interface{}
	c.ShouldBindJSON(&data)

	id, err := h.Service.Create(resource, data, claims.RoleID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Created"})
}

func (h *Handler) Update(c *gin.Context) {
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)

	resource := c.Param("resource")
	id := c.Param("id")
	var data map[string]interface{}
	c.ShouldBindJSON(&data)

	err := h.Service.Update(resource, id, data, claims.RoleID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	user, _ := c.Get("user")
	claims := user.(*utils.Claims)

	resource := c.Param("resource")
	id := c.Param("id")

	err := h.Service.Delete(resource, id, claims.RoleID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
