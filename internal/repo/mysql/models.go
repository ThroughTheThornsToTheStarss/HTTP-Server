package mysql

import (
	"time"
)

type Account struct {
	ID           string `gorm:"primaryKey;size:64"`
	AccessToken  string
	RefreshToken string
	Expires      int64

	Integrations []Integration `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
	Contacts     []Contact     `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Integration struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement"`
	AccountID           string `gorm:"index;size:64;not null"`
	SecretKey           string
	ClientID            string
	RedirectURL         string
	AuthenticationCode  string
	AccessToken         string
	RefreshToken        string
	AccessTokenExpires  int64
	RefreshTokenExpires int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Contact struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	AccountID string `gorm:"index;size:64;not null"`

	Name  string
	Email *string

	CreatedAt time.Time
	UpdatedAt time.Time
}
