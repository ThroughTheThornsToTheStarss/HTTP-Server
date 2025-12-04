package repo

import "git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"

type Repository interface {
	CreateAccount(acc *domain.Account) error
	GetAllAccounts() ([]*domain.Account, error)
	
	DeleteAccount(accountID string) error
	UpdateAccount(acc *domain.Account) error

	CreateIntegration(in *domain.Integration) error
	GetIntegrationsByAccountID(accountID string) ([]*domain.Integration, error)
}
