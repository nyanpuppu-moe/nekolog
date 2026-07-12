package service

import (
	"errors"

	"nekolog/internal/dto"
	"nekolog/internal/model"
	"nekolog/internal/repository"
	"nekolog/internal/utils"
)

type ArticleService struct {
	articleRepo *repository.ArticleRepository
	userRepo    *repository.UserRepository
}

func NewArticleService(
	articleRepo *repository.ArticleRepository,
	userRepo *repository.UserRepository,
) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

func (s *ArticleService) Get(username string, title string) (*model.Article, error) {
	user, err := s.userRepo.FindByName(username)
	if err != nil {
		return nil, err
	}

	return s.articleRepo.FindByAuthorIDAndTitle(user.ID, title)
}

func (s *ArticleService) Post(userID model.UserID, req dto.ArticlePostRequest) error {
	article, _ := s.articleRepo.FindByAuthorIDAndTitle(userID, req.Title)
	if article != nil {
		return errors.New("Already same title article")
	}

	newArticle := &model.Article{
		AuthorID: userID,
		Title:    req.Title,
	}

	err := s.articleRepo.Create(newArticle, req.Content)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArticleService) Patch(userID model.UserID, title string, req dto.ArticlePatchRequest) error {
	currentArticle, err := s.articleRepo.FindByAuthorIDAndTitle(userID, title)
	if err != nil {
		return err
	}

	updates := make(utils.Object)
	updates["content"] = req.Content

	return s.articleRepo.Update(currentArticle.ID, updates)
}

func (s *ArticleService) Delete(userID model.UserID, title string) error {
	return s.articleRepo.Delete(userID, title)
}
