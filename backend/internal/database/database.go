package database

import (
	"nekolog/internal/log"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

func Connect(dbPath string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: log.NewGormLogger(),
	}

	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}
