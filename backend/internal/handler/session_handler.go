package handler

import (
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/log"
	"nekolog/internal/service"
	"nekolog/internal/web"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(s *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: s}
}

func (h *SessionHandler) Register(c *web.Context) {
	var req dto.SessionRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("올바르지 않은 요청입니다: %v", err)
		c.JSON(
			http.StatusBadRequest,
			web.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	err := h.sessionService.Register(req)
	if err != nil {
		log.Warn("사용자 등록에 실패하였습니다: %v", err)
		c.JSON(
			http.StatusConflict,
			web.Object{
				"error": "사용자 등록에 실패하였습니다",
			},
		)
		return
	}

	log.Info("사용자를 등록하였습니다: %s", req.Name)
	c.JSON(
		http.StatusCreated,
		web.Object{
			"message": "사용자를 등록하였습니다",
		},
	)
}

func (h *SessionHandler) Login(c *web.Context) {
	var req dto.SessionLoginRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("올바르지 않은 요청입니다: %v", err)
		c.JSON(
			http.StatusBadRequest,
			web.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	user, err := h.sessionService.Login(req)
	if err != nil {
		log.Warn("사용자 이름 또는 비밀번호가 올바르지 않습니다")
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "사용자 이름 또는 비밀번호가 올바르지 않습니다",
			},
		)
		return
	}

	c.SessionSet("user_id", user.ID)
	if err := c.SessionSave(); err != nil {
		log.Warn("세션을 저장하지 못했습니다: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			web.Object{
				"error": "세션을 저장하지 못했습니다",
			},
		)
		return
	}

	log.Info("로그인에 성공하였습니다: %s", user.Name)
	c.JSON(
		http.StatusOK,
		web.Object{
			"message": "로그인에 성공하였습니다",
			"user":    user,
		},
	)
}

func (h *SessionHandler) Logout(c *web.Context) {
	c.SessionClear()

	if err := c.SessionSave(); err != nil {
		log.Warn("로그아웃에 실패하였습니다: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			web.Object{
				"error": "로그아웃에 실패하였습니다",
			},
		)
		return
	}

	log.Info("로그아웃에 성공하였습니다")
	c.JSON(
		http.StatusOK,
		web.Object{
			"message": "로그아웃에 성공하였습니다",
		},
	)
}
