package storage

import (
	"os"
	"path/filepath"

	"nekolog/internal/model"
)

type ContentStorage struct {
	storagePath string
}

func NewContentStorage(storagePath string) *ContentStorage {
	return &ContentStorage{
		storagePath: storagePath,
	}
}

func (s *ContentStorage) Create(article *model.Article, content string) error {
	fullPath := filepath.Join(
		s.storagePath,
		toStringUserID(article.AuthorID),
		article.Title,
	)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	data := []byte(content)

	return os.WriteFile(fullPath, data, 0644)
}

func (s *ContentStorage) Delete(authorID model.UserID, title string) error {
	fullPath := filepath.Join(
		s.storagePath,
		toStringUserID(authorID),
		title,
	)

	return os.Remove(fullPath)
}
