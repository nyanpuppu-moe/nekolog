package handler

import (
	"io"
	"net/http"

	"nekolog/internal/dto"
	"nekolog/internal/log"
	"nekolog/internal/model"
	"nekolog/internal/service"
	"nekolog/internal/utils"
	"nekolog/internal/web"
)

type AssetHandler struct {
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

func (h *AssetHandler) Get(c *web.Context) {
	id, err := utils.StringToUint[model.AssetID](c.Param("id"))
	if err != nil {
		log.Warn("올바르지 않은 아이디입니다: %d", id)
		c.JSON(
			http.StatusBadRequest,
			web.Object{
				"error": "올바르지 않은 아이디입니다",
			},
		)
		return
	}

	asset, err := h.assetService.Get(id)
	if err != nil {
		log.Warn("존재하지 않는 에셋입니다: %d", id)
		c.JSON(
			http.StatusNotFound,
			web.Object{
				"error": "존재하지 않는 에셋입니다",
			},
		)
		return
	}

	log.Info("에셋이 조회되었습니다: %d", id)
	c.JSON(
		http.StatusOK,
		web.Object{
			"message": "에셋이 조회되었습니다",
			"user":    asset,
		},
	)
}

func (h *AssetHandler) Post(c *web.Context) {
	sessionUserID := c.SessionGet("user_id")
	if sessionUserID == nil {
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
		c.JSON(
			http.StatusUnauthorized,
			web.Object{
				"error": "올바르지 않은 유저입니다",
			},
		)
		return
	}

	var req dto.AssetPostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			web.Object{
				"error": "올바르지 않은 요청입니다",
			},
		)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			web.Object{
				"error": "파일이 존재하지 않습니다",
			},
		)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			web.Object{
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
			web.Object{
				"error": "파일을 읽을 수 없습니다",
			},
		)
		return
	}

	if err := h.assetService.Post(userID, data, req); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			web.Object{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusCreated, web.Object{
		"message": "파일이 생성되었습니다",
	})
}
