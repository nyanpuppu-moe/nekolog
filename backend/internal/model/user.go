package model

type UserID uint64

type User struct {
	ID          UserID `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null;unique"`
	DisplayName string
	Password    string `gorm:"not null"`
}
