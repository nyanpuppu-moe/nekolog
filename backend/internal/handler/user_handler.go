package handler

import (
	"net/http"

	"nekolog/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) Get(c *gin.Context) {
	name := c.Param("username")

	user, err := h.userService.Get(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get user successfully",
		"user":    user,
	})
}
