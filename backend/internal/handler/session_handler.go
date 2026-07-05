package handler

import (
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(s *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: s}
}

func (h *SessionHandler) Register(c *gin.Context) {
	var req dto.SessionRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := h.sessionService.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Faild to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Register successfully"})
}

func (h *SessionHandler) Login(c *gin.Context) {
	var req dto.SessionLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user, err := h.sessionService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faild to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"user":    user,
	})
}

func (h *SessionHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
