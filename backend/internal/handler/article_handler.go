package handler

import (
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/engine"
	"nekolog/internal/model"
	"nekolog/internal/service"
)

type ArticleHandler struct {
	articleService *service.ArticleService
	userService    *service.UserService
}

func NewArticleHandler(
	articleService *service.ArticleService,
	userService *service.UserService,
) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		userService:    userService,
	}
}

func (h *ArticleHandler) Get(c *engine.Context) {
	username := c.Param("username")
	title := c.Param("title")

	article, err := h.articleService.Get(username, title)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			engine.Object{
				"error": "존재하지 않는 아티클입니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		engine.Object{
			"message": "아티클이 조회되었습니다",
			"article": article,
		},
	)
}

func (h *ArticleHandler) Post(c *engine.Context) {
	sessionUserID := c.SessionGet("user_id")
	if sessionUserID == nil {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	var req dto.ArticlePostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			engine.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	if err := h.articleService.Post(userID, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		engine.Object{
			"message": "아티클이 생성되었습니다",
		},
	)
}

func (h *ArticleHandler) Patch(c *engine.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID := c.SessionGet("user_id")
	if sessionUserID == nil {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	user, err := h.userService.FindUserByID(userID)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	if user.Name != username {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	var req dto.ArticlePatchRequest
	if err = c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			engine.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	if err = h.articleService.Patch(userID, title, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		engine.Object{
			"message": "아티클이 수정되었습니다",
		},
	)
}

func (h *ArticleHandler) Delete(c *engine.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID := c.SessionGet("user_id")
	if sessionUserID == nil {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	user, err := h.userService.FindUserByID(userID)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	if user.Name != username {
		c.JSON(
			http.StatusUnauthorized,
			engine.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	if err = h.articleService.Delete(userID, title); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, engine.Object{
		"message": "아티클이 삭제되었습니다",
	})
}
