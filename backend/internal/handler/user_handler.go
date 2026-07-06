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
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "존재하지 않는 사용자입니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "사용자를 조회하였습니다",
			"user":    user,
		},
	)
}
