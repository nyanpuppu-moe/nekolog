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
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	err := h.sessionService.Register(req)
	if err != nil {
		c.JSON(
			http.StatusConflict,
			gin.H{
				"error": "사용자 등록에 실패하였습니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "사용자를 등록하였습니다",
		},
	)
}

func (h *SessionHandler) Login(c *gin.Context) {
	var req dto.SessionLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	user, err := h.sessionService.Login(req)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "사용자 이름 또는 비밀번호가 올바르지 않습니다",
			},
		)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "세션을 저장하지 못했습니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "로그인에 성공하였습니다",
			"user":    user,
		},
	)
}

func (h *SessionHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()

	if err := session.Save(); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "로그아웃에 실패하였습니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "로그아웃에 성공하였습니다",
		},
	)
}
