package handler

import (
	"io"
	"net/http"
	"strconv"

	"nekolog/internal/dto"
	"nekolog/internal/model"
	"nekolog/internal/service"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (h *AssetHandler) Get(c *gin.Context) {
	idParam := c.Param("id")

	idU64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id",
		})
		return
	}

	id := model.AssetID(idU64)

	asset, err := h.assetService.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found asset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get user successfully",
		"user":    asset,
	})
}

func (h *AssetHandler) Post(c *gin.Context) {
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

	var req dto.AssetPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.assetService.Post(userID, data, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post asset successfully",
	})
}
