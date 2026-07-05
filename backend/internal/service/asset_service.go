package service

import (
	"nekolog/internal/dto"
	"nekolog/internal/model"
	"nekolog/internal/repository"
)

type AssetService struct {
	assetRepo *repository.AssetRepository
	userRepo  *repository.UserRepository
}

func NewAssetService(
	assetRepo *repository.AssetRepository,
	userRepo *repository.UserRepository,
) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
		userRepo:  userRepo,
	}
}

func (s *AssetService) Get(id model.AssetID) (*model.Asset, error) {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *AssetService) Post(user *model.User, data []byte, req dto.AssetPostRequest) error {
	newAsset := &model.Asset{
		AuthorID: user.ID,
		Author:   *user,
		Filename: req.Filename,
	}

	return s.assetRepo.Create(newAsset, data)
}
