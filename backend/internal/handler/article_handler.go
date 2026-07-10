package handler

import (
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/log"
	"nekolog/internal/model"
	"nekolog/internal/service"
	"nekolog/internal/web"
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

func (h *ArticleHandler) Get(c *web.Context) {
	username := c.Param("username")
	title := c.Param("title")

	article, err := h.articleService.Get(username, title)
	if err != nil {
		log.Warn("존재하지 않는 아티클입니다: %s/%s", username, title)
		c.JSON(
			http.StatusNotFound,
			web.Object{
				"error": "존재하지 않는 아티클입니다",
			},
		)
		return
	}

	log.Info("아티클이 조회되었습니다: %s/%s", username, title)
	c.JSON(
		http.StatusOK,
		web.Object{
			"message": "아티클이 조회되었습니다",
			"article": article,
		},
	)
}

func (h *ArticleHandler) Post(c *web.Context) {
	sessionUserID := c.SessionGet("user_id")
	if sessionUserID == nil {
		log.Warn("인증되지 않은 유저입니다")
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		log.Warn("올바르지 않은 유저입니다: %d", userID)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	user, err := h.userService.FindUserByID(userID)
	if err != nil {
		log.Warn("올바르지 않은 유저입니다: %s", user.Name)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	var req dto.ArticlePostRequest
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

	if err := h.articleService.Post(userID, req); err != nil {
		log.Warn("아티클 생성에 실패하였습니다: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			web.Object{
				"error": "아티클 생성에 실패하였습니다",
			},
		)
		return
	}

	log.Info("아티클이 생성되었습니다: %s/%s", user.Name, req.Title)
	c.JSON(
		http.StatusOK,
		web.Object{
			"message": "아티클이 생성되었습니다",
		},
	)
}

func (h *ArticleHandler) Patch(c *web.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID := c.SessionGet("user_id")
	if sessionUserID == nil {
		log.Warn("인증되지 않은 유저입니다")
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		log.Warn("올바르지 않은 유저입니다: %d", userID)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	user, err := h.userService.FindUserByID(userID)
	if err != nil {
		log.Warn("올바르지 않은 유저입니다: %s", user.Name)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	if user.Name != username {
		log.Warn("인증되지 않은 유저입니다 %v", username)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	var req dto.ArticlePatchRequest
	if err = c.BindJSON(&req); err != nil {
		log.Warn("올바르지 않은 요청입니다: %v", err)
		c.JSON(
			http.StatusBadRequest,
			web.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	if err = h.articleService.Patch(userID, title, req); err != nil {
		log.Warn("아티클 수정에 실패하였습니다: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			web.Object{
				"error": "아티클 수정에 실패하였습니다",
			},
		)
		return
	}

	log.Info("아티클이 수정되었습니다: %s/%s", username, title)
	c.JSON(
		http.StatusOK,
		web.Object{
			"message": "아티클이 수정되었습니다",
		},
	)
}

func (h *ArticleHandler) Delete(c *web.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID := c.SessionGet("user_id")
	if sessionUserID == nil {
		log.Warn("인증되지 않은 유저입니다")
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		log.Warn("올바르지 않은 유저입니다: %d", userID)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	user, err := h.userService.FindUserByID(userID)
	if err != nil {
		log.Warn("올바르지 않은 유저입니다: %s", user.Name)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	if user.Name != username {
		log.Warn("인증되지 않은 유저입니다 %v", username)
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	if err = h.articleService.Delete(userID, title); err != nil {
		log.Warn("아티클 삭제에 실패하였습니다: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			web.Object{
				"error": "아티클 삭제에 실패하였습니다",
			},
		)
		return
	}

	log.Info("아티클이 삭제되었습니다: %s/%s", username, title)
	c.JSON(http.StatusOK, web.Object{
		"message": "아티클이 삭제되었습니다",
	})
}
