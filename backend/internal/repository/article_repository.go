package repository

import (
	"nekolog/internal/model"
	"nekolog/internal/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	db             *gorm.DB
	contentStorage *storage.ContentStorage
}

func NewArticleRepository(
	db *gorm.DB,
	contentStorage *storage.ContentStorage,
) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(article *model.Article, content string) error {
	if err := r.contentStorage.Create(article.AuthorID, article.Title, content); err != nil {
		return err
	}

	return r.db.
		Create(article).
		Error
}

func (r *ArticleRepository) FindByAuthorIDAndTitle(authorID model.UserID, title string) (*model.Article, error) {
	var article model.Article

	err := r.db.
		Where(&model.Article{
			AuthorID: authorID,
			Title:    title,
		}).
		First(&article).
		Error

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *ArticleRepository) Update(id model.ArticleID, updates gin.H) error {
	return r.db.
		Model(&model.Article{}).
		Where(&model.Article{
			ID: id,
		}).
		Updates(updates).
		Error
}

func (r *ArticleRepository) Delete(authorID model.UserID, title string) error {
	if err := r.contentStorage.Delete(authorID, title); err != nil {
		return err
	}

	return r.db.
		Delete(&model.Article{
			AuthorID: authorID,
			Title:    title,
		}).
		Error
}
