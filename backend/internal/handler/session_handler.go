package handler

import (
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/engine"
	"nekolog/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(s *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: s}
}

func (h *SessionHandler) Register(c *engine.Context) {
	var req dto.SessionRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			engine.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	err := h.sessionService.Register(req)
	if err != nil {
		c.JSON(
			http.StatusConflict,
			engine.Object{
				"error": "사용자 등록에 실패하였습니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		engine.Object{
			"message": "사용자를 등록하였습니다",
		},
	)
}

func (h *SessionHandler) Login(c *engine.Context) {
	var req dto.SessionLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			engine.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	user, err := h.sessionService.Login(req)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "사용자 이름 또는 비밀번호가 올바르지 않습니다",
			},
		)
		return
	}

	c.SessionSet("user_id", user.ID)
	if err := c.SessionSave(); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": "세션을 저장하지 못했습니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		engine.Object{
			"message": "로그인에 성공하였습니다",
			"user":    user,
		},
	)
}

func (h *SessionHandler) Logout(c *engine.Context) {
	c.SessionClear()

	if err := c.SessionSave(); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": "로그아웃에 실패하였습니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		engine.Object{
			"message": "로그아웃에 성공하였습니다",
		},
	)
}
