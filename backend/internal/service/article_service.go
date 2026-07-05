// TODO: Add error message

package service

import (
	"errors"
	"nekolog/internal/dto"
	"nekolog/internal/model"
	"nekolog/internal/repository"

	"github.com/gin-gonic/gin"
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

func (s *ArticleService) Post(user *model.User, req dto.ArticlePostRequest) error {
	_, err := s.articleRepo.FindByAuthorIDAndTitle(user.ID, req.Title)
	if err == nil {
		return errors.New("Already same title article")
	}

	newArticle := &model.Article{
		AuthorID: user.ID,
		Author:   *user,
		Title:    req.Title,
		Content:  req.Content,
	}

	if err := s.articleRepo.Create(newArticle); err != nil {
		return err
	}

	return nil
}

func (s *ArticleService) Patch(user *model.User, title string, req dto.ArticlePatchRequest) error {
	currentArticle, err := s.articleRepo.FindByAuthorIDAndTitle(user.ID, title)
	if err != nil {
		return err
	}

	updates := make(gin.H)
	updates["content"] = req.Content

	return s.articleRepo.Update(currentArticle.ID, updates)
}

func (s *ArticleService) Delete(user *model.User, title string) error {
	article, err := s.articleRepo.FindByAuthorIDAndTitle(user.ID, title)
	if err != nil {
		return err
	}

	return s.articleRepo.Delete(article.ID)
}
