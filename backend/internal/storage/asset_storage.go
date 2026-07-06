package storage

import (
	"os"
	"path/filepath"

	"nekolog/internal/model"
	"nekolog/pkg/utils"
)

type AssetStorage struct {
	storagePath string
}

func NewAssetStorage(storagePath string) *AssetStorage {
	return &AssetStorage{
		storagePath: storagePath,
	}
}

func (s *AssetStorage) Create(authorID model.UserID, assetID model.AssetID, data []byte) error {
	fullPath := filepath.Join(
		s.storagePath,
		utils.ToString(authorID),
		utils.ToString(assetID),
	)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}

func (s *AssetStorage) Delete(authorID model.UserID, assetID model.AssetID) error {
	fullPath := filepath.Join(
		s.storagePath,
		utils.ToString(authorID),
		utils.ToString(assetID),
	)

	return os.Remove(fullPath)
}
