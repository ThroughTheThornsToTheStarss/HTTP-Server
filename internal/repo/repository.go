package repo

import "git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"

type Repository interface {
	CreateAccount(acc *domain.Account) error
	GetAllAccounts() ([]*domain.Account, error)
	GetAccountByID(accountID uint64) (*domain.Account, error)
	DeleteAccount(accountID uint64) error
	UpdateAccount(acc *domain.Account) error

	CreateIntegration(in *domain.Integration) error
	GetIntegrationsByAccountID(accountID uint64) ([]*domain.Integration, error)

	SaveContacts(accountID uint64, contacts []*domain.Contact) error
	GetContactsByAccountID(accountID uint64) ([]*domain.Contact, error)

	UpsertContactFromWebhook(accountID uint64, amoID int64, name string, email *string, status string) (contactID uint, err error)
	FindContactIDByAmoID(accountID uint64, amoID int64) (uint, bool, error)

	AddSyncHistory(contactID uint, status string, message string) error
	TrimSyncHistory(contactID uint, keepLast int) error
}
