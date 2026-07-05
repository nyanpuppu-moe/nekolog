package storage

import (
	"strconv"

	"os"
	"path/filepath"

	"nekolog/internal/model"
)

type AssetStorage struct {
	storagePath string
}

func toStringUserID(id model.UserID) string {
	return strconv.FormatUint(uint64(id), 10)
}

func toStringAssetID(id model.AssetID) string {
	return strconv.FormatUint(uint64(id), 10)
}

func NewAssetStorage(storagePath string) *AssetStorage {
	return &AssetStorage{
		storagePath: storagePath,
	}
}

func (s *AssetStorage) Create(asset *model.Asset, data []byte) error {
	fullPath := filepath.Join(
		s.storagePath,
		toStringUserID(asset.AuthorID),
		toStringAssetID(asset.ID),
	)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}

func (s *AssetStorage) Delete(authorID model.UserID, assetID model.AssetID) error {
	fullPath := filepath.Join(
		s.storagePath,
		toStringUserID(authorID),
		toStringAssetID(assetID),
	)

	return os.Remove(fullPath)
}
