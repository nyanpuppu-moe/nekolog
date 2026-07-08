package handler

import (
	"io"
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/engine"
	"nekolog/internal/model"
	"nekolog/internal/service"
	"nekolog/internal/utils"
)

type AssetHandler struct {
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

func (h *AssetHandler) Get(c *engine.Context) {
	id, err := utils.StringToUint[model.AssetID](c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			engine.Object{
				"error": "올바르지 않은 아이디입니다",
			},
		)
		return
	}

	asset, err := h.assetService.Get(id)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			engine.Object{
				"error": "존재하지 않는 에셋입니다",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		engine.Object{
			"message": "에셋이 조회되었습니다",
			"user":    asset,
		},
	)
}

func (h *AssetHandler) Post(c *engine.Context) {
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

	var req dto.AssetPostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			engine.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			engine.Object{
				"error": "파일이 존재하지 않습니다",
			},
		)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": "파일을 열 수 없습니다",
			},
		)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": "파일을 읽을 수 없습니다",
			},
		)
		return
	}

	if err := h.assetService.Post(userID, data, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			engine.Object{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusCreated, engine.Object{
		"message": "파일이 생성되었습니다",
	})
}
