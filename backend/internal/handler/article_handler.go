package handler

import (
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/model"
	"nekolog/internal/service"

	"github.com/gin-gonic/gin"
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

func (h *ArticleHandler) Get(c *gin.Context) {
	username := c.Param("username")
	title := c.Param("title")

	article, err := h.articleService.Get(username, title)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "존재하지 않는 아티클입니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "아티클이 조회되었습니다",
			"article": article,
		},
	)
}

func (h *ArticleHandler) Post(c *gin.Context) {
	sessionUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	var req dto.ArticlePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	if err := h.articleService.Post(userID, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "아티클이 생성되었습니다",
		},
	)
}

func (h *ArticleHandler) Patch(c *gin.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	user, err := h.userService.FindUserByID(userID)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	if user.Name != username {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	var req dto.ArticlePatchRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	if err = h.articleService.Patch(userID, title, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "아티클이 수정되었습니다",
		},
	)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	user, err := h.userService.FindUserByID(userID)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	if user.Name != username {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "인증되지 않은 유저입니다",
			},
		)
		return
	}

	if err = h.articleService.Delete(userID, title); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "아티클이 삭제되었습니다",
	})
}
