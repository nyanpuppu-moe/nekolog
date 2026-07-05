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

func NewArticleHandler(s *service.ArticleService, us *service.UserService) *ArticleHandler {
	return &ArticleHandler{articleService: s, userService: us}
}

func (h *ArticleHandler) Get(c *gin.Context) {
	username := c.Param("username")
	title := c.Param("title")

	article, err := h.articleService.Get(username, title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get article successfully",
		"article": article,
	})
}

func (h *ArticleHandler) Post(c *gin.Context) {
	sessionUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not verified user"})
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
		return
	}

	var req dto.ArticlePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.articleService.Post(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post article successfully",
	})
}

func (h *ArticleHandler) Patch(c *gin.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not verified user"})
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if user.Name != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not verified user"})
		return
	}

	var req dto.ArticlePatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err = h.articleService.Patch(userID, title, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patch article successfully",
	})
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	username := c.Param("username")
	title := c.Param("title")

	sessionUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not verified user"})
		return
	}

	userID, ok := sessionUserID.(model.UserID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if user.Name != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not verified user"})
		return
	}

	if err = h.articleService.Delete(userID, title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete article successfully",
	})
}
