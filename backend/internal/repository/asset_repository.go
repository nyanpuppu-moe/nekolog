package repository

import (
	"nekolog/internal/model"
	"nekolog/internal/storage"

	"gorm.io/gorm"
)

type AssetRepository struct {
	db      *gorm.DB
	storage *storage.AssetStorage
}

func NewAssetRepository(
	db *gorm.DB,
	storage *storage.AssetStorage,
) *AssetRepository {
	return &AssetRepository{
		db:      db,
		storage: storage,
	}
}

func (r *AssetRepository) Create(asset *model.Asset, data []byte) error {
	err := r.storage.Create(asset, data)
	if err != nil {
		return err
	}

	return r.db.
		Create(asset).
		Error
}

func (r *AssetRepository) Delete(userID model.UserID, assetID model.AssetID) error {
	err := r.storage.Delete(userID, assetID)
	if err != nil {
		return err
	}

	return r.db.
		Delete(&model.Asset{
			ID: assetID,
		}).
		Error
}

func (r *AssetRepository) FindByID(id model.AssetID) (*model.Asset, error) {
	var asset model.Asset

	err := r.db.
		Where(&model.Asset{
			ID: id,
		}).
		First(&asset).
		Error

	if err != nil {
		return nil, err
	}

	return &asset, nil
}
