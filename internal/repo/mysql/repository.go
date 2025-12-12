package mysql

import (
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo"
	"gorm.io/gorm"
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
	return r.db.Create(&model).Error
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

func (r *GormRepository) DeleteAccount(accountID string) error {
	return r.db.Where("id = ?", accountID).Delete(&Account{}).Error
}

func (r *GormRepository) UpdateAccount(acc *domain.Account) error {
	return r.db.Model(&Account{}).
		Where("id = ?", acc.ID).
		Updates(map[string]any{
			"access_token":  acc.AccessToken,
			"refresh_token": acc.RefreshToken,
			"expires":       acc.Expires,
		}).Error
}

func (r *GormRepository) CreateIntegration(in *domain.Integration) error {
	model := integrationToModel(in)
	return r.db.Create(&model).Error
}

func (r *GormRepository) GetIntegrationsByAccountID(accountID string) ([]*domain.Integration, error) {
	var models []Integration
	if err := r.db.
		Where("account_id = ?", accountID).
		Find(&models).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.Integration, 0, len(models))
	for i := range models {
		res = append(res, integrationFromModel(&models[i]))
	}
	return res, nil
}

func accountToModel(a *domain.Account) Account {
	if a == nil {
		return Account{}
	}
	return Account{
		ID:           a.ID,
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		Expires:      a.Expires,
	}
}

func accountFromModel(m *Account) *domain.Account {
	if m == nil {
		return nil
	}
	return &domain.Account{
		ID:           m.ID,
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		Expires:      m.Expires,
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
	}
}

func (r *GormRepository) SaveContacts(accountID string, contacts []*domain.Contact) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id = ?", accountID).Delete(&Contact{}).Error; err != nil {
			return err
		}

		if len(contacts) == 0 {
			return nil
		}

		models := make([]Contact, 0, len(contacts))
		for _, c := range contacts {
			models = append(models, Contact{
				AccountID: accountID,
				Name:      c.Name,
				Email:     c.Email,
			})
		}

		return tx.Create(&models).Error
	})
}

func (r *GormRepository) GetContactsByAccountID(accountID string) ([]*domain.Contact, error) {
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

func contactFromModel(m *Contact) *domain.Contact {
	if m == nil {
		return nil
	}
	return &domain.Contact{
		ID:        m.ID,
		AccountID: m.AccountID,
		Name:      m.Name,
		Email:     m.Email,
	}
}
