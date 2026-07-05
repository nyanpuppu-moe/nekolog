package repository

import (
	"nekolog/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(article *model.Article) error {
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

func (r *ArticleRepository) Delete(id model.ArticleID) error {
	return r.db.
		Delete(&model.Article{
			ID: id,
		}).
		Error
}
