package mysql

import "time"

type Account struct {
	ID uint64 `gorm:"primaryKey;autoIncrement:false"`

	Referer string `gorm:"size:255;not null"`

	AccessToken  string `gorm:"type:text"`
	RefreshToken string `gorm:"type:text"`
	TokenType    string `gorm:"size:32"`
	ExpiresIn    int64

	Integrations []Integration `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
	Contacts     []Contact     `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool `gorm:"not null;default:true"`
}

type Integration struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement"`
	AccountID          uint64 `gorm:"not null;uniqueIndex"`
	SecretKey          string
	ClientID           string
	RedirectURL        string
	AuthenticationCode string

	CreatedAt    time.Time
	UpdatedAt    time.Time
	UnisenderKey string `gorm:"type:text"`
}

type Contact struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	AccountID uint64 `gorm:"index;not null"`

	Name  string
	Email *string

	CreatedAt time.Time
	UpdatedAt time.Time
}
