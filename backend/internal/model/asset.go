package model

import "time"

type AssetID uint64

type Asset struct {
	ID        AssetID   `gorm:"primaryKey;autoIncrement"`
	AuthorID  UserID    `gorm:"not null;uniqueIndex:uniq_author_title"`
	Author    User      `gorm:"constraint:OnDelete:CASCADE;"`
	Filename  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
