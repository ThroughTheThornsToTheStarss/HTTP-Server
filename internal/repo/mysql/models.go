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
	AccountID uint64 `gorm:"not null;uniqueIndex:uniq_acc_email"`

	Name  string  `gorm:"size:255"`
	Email *string `gorm:"size:255;uniqueIndex:uniq_acc_email"`

	Status    string `gorm:"size:16;not null;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type SyncHistory struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ContactID uint   `gorm:"index;not null"`
	Status    string `gorm:"size:16;not null"`
	Message   string `gorm:"type:text"`

	CreatedAt time.Time
}
