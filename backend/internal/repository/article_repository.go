package repository

import (
	"nekolog/internal/model"
	"nekolog/internal/storage"
	"nekolog/internal/web"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	dataBase       *gorm.DB
	contentStorage *storage.ContentStorage
}

func NewArticleRepository(
	dataBase *gorm.DB,
	contentStorage *storage.ContentStorage,
) *ArticleRepository {
	return &ArticleRepository{
		dataBase:       dataBase,
		contentStorage: contentStorage,
	}
}

func (r *ArticleRepository) Create(article *model.Article, content string) error {
	if err := r.contentStorage.Create(article.AuthorID, article.Title, content); err != nil {
		return err
	}

	return r.dataBase.
		Create(article).
		Error
}

func (r *ArticleRepository) FindByAuthorIDAndTitle(authorID model.UserID, title string) (*model.Article, error) {
	var article model.Article

	err := r.dataBase.
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

func (r *ArticleRepository) Update(id model.ArticleID, updates web.Object) error {
	return r.dataBase.
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

	return r.dataBase.
		Delete(&model.Article{
			AuthorID: authorID,
			Title:    title,
		}).
		Error
}
