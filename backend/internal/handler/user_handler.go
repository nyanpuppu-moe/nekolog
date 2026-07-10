package handler

import (
	"net/http"

	"nekolog/internal/log"
	"nekolog/internal/service"
	"nekolog/internal/web"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) Get(c *web.Context) {
	name := c.Param("username")

	user, err := h.userService.Get(name)
	if err != nil {
		log.Warn("존재하지 않는 사용자입니다: %s", name)
		c.JSON(
			http.StatusNotFound,
			web.Object{
				"error": "존재하지 않는 사용자입니다",
			},
		)
		return
	}

	log.Info("사용자를 조회하였습니다: %s", name)
	c.JSON(
		http.StatusOK,
		web.Object{
			"message": "사용자를 조회하였습니다",
			"user":    user,
		},
	)
}
