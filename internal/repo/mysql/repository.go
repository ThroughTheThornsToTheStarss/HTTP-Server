package mysql

import (
	"errors"
	"strings"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ repo.Repository = (*GormRepository)(nil)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateAccount(acc *domain.Account) error {
	model := accountToModel(acc)
	model.IsActive = true
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"referer",
			"access_token",
			"refresh_token",
			"token_type",
			"expires_in",
			"is_active",
		}),
	}).Create(&model).Error
}

func (r *GormRepository) GetAllAccounts() ([]*domain.Account, error) {
	var models []Account
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.Account, 0, len(models))
	for i := range models {
		res = append(res, accountFromModel(&models[i]))
	}
	return res, nil
}

func (r *GormRepository) GetAccountByID(accountID uint64) (*domain.Account, error) {
	var m Account
	if err := r.db.First(&m, "id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return accountFromModel(&m), nil
}

func (r *GormRepository) DeleteAccount(accountID uint64) error {
	return r.db.Model(&Account{}).
		Where("id = ?", accountID).
		Update("is_active", false).Error
}

func (r *GormRepository) UpdateAccount(acc *domain.Account) error {
	if acc == nil {
		return errors.New("nil account")
	}

	return r.db.Model(&Account{}).
		Where("id = ?", acc.ID).
		Updates(map[string]any{
			"referer":       acc.Referer,
			"access_token":  acc.AccessToken,
			"refresh_token": acc.RefreshToken,
			"token_type":    acc.TokenType,
			"expires_in":    acc.ExpiresIn,
		}).Error
}

func (r *GormRepository) CreateIntegration(in *domain.Integration) error {
	model := integrationToModel(in)

	var existing Integration
	err := r.db.Where("account_id = ?", model.AccountID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(&model).Error
	}
	if err != nil {
		return err
	}

	updates := map[string]any{}
	if model.UnisenderKey != "" {
		updates["unisender_key"] = model.UnisenderKey
	}
	if model.SecretKey != "" {
		updates["secret_key"] = model.SecretKey
	}
	if model.ClientID != "" {
		updates["client_id"] = model.ClientID
	}
	if model.RedirectURL != "" {
		updates["redirect_url"] = model.RedirectURL
	}
	if model.AuthenticationCode != "" {
		updates["authentication_code"] = model.AuthenticationCode
	}

	return r.db.Model(&existing).Updates(updates).Error
}

func (r *GormRepository) GetIntegrationsByAccountID(accountID uint64) ([]*domain.Integration, error) {
	var models []Integration
	if err := r.db.Where("account_id = ?", accountID).Find(&models).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.Integration, 0, len(models))
	for i := range models {
		res = append(res, integrationFromModel(&models[i]))
	}
	return res, nil
}

func (r *GormRepository) SaveContacts(accountID uint64, contacts []*domain.Contact) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Contact{}).
			Where("account_id = ?", accountID).
			Update("status", "deleted").Error; err != nil {
			return err
		}

		if len(contacts) == 0 {
			return nil
		}

		models := make([]Contact, 0, len(contacts))
		seen := map[string]struct{}{}

		for _, c := range contacts {
			if c == nil || c.Email == nil {
				continue
			}
			email := strings.TrimSpace(*c.Email)
			if email == "" {
				continue
			}
			if _, ok := seen[email]; ok {
				continue
			}
			seen[email] = struct{}{}

			emailCopy := email
			models = append(models, Contact{
				AccountID: accountID,
				Name:      c.Name,
				Email:     &emailCopy,
				Status:    "active",
			})
		}

		if len(models) == 0 {
			return nil
		}

		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "account_id"}, {Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"status",
				"updated_at",
			}),
		}).Create(&models).Error
	})
}

func (r *GormRepository) GetContactsByAccountID(accountID uint64) ([]*domain.Contact, error) {
	var models []Contact
	if err := r.db.Where("account_id = ?", accountID).Find(&models).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.Contact, 0, len(models))
	for _, m := range models {
		res = append(res, contactFromModel(&m))
	}
	return res, nil
}

func accountToModel(a *domain.Account) Account {
	if a == nil {
		return Account{}
	}
	return Account{
		ID:           a.ID,
		Referer:      a.Referer,
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		TokenType:    a.TokenType,
		ExpiresIn:    a.ExpiresIn,
		IsActive:     a.IsActive,
	}
}

func accountFromModel(m *Account) *domain.Account {
	if m == nil {
		return nil
	}
	return &domain.Account{
		ID:           m.ID,
		Referer:      m.Referer,
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		TokenType:    m.TokenType,
		ExpiresIn:    m.ExpiresIn,
		IsActive:     m.IsActive,
	}
}

func integrationToModel(in *domain.Integration) Integration {
	if in == nil {
		return Integration{}
	}
	return Integration{
		AccountID:          in.AccountID,
		SecretKey:          in.SecretKey,
		ClientID:           in.ClientID,
		RedirectURL:        in.RedirectURL,
		AuthenticationCode: in.AuthenticationCode,
		UnisenderKey:       in.UnisenderKey,
	}
}

func integrationFromModel(m *Integration) *domain.Integration {
	if m == nil {
		return nil
	}
	return &domain.Integration{
		AccountID:          m.AccountID,
		SecretKey:          m.SecretKey,
		ClientID:           m.ClientID,
		RedirectURL:        m.RedirectURL,
		AuthenticationCode: m.AuthenticationCode,
		UnisenderKey:       m.UnisenderKey,
	}
}

func contactFromModel(m *Contact) *domain.Contact {
	if m == nil {
		return nil
	}
	return &domain.Contact{
		ID:        m.ID,
		AccountID: m.AccountID,
		Name:      m.Name,
		Email:     m.Email,
		Status:    m.Status,
	}
}
