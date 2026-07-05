package model

type ArticleID uint64

type Article struct {
	ID       ArticleID `gorm:"primaryKey;autoIncrement"`
	AuthorID UserID    `gorm:"not null;uniqueIndex:uniq_author_title"`
	Author   User      `gorm:"constraint:OnDelete:CASCADE;"`
	Title    string    `gorm:"not null;uniqueIndex:uniq_author_title"`
	Content  string    `gorm:"not null;type:text"`
	View     uint      `gorm:"default:0"`
}
