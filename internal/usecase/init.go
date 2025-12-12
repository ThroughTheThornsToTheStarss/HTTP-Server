package usecase

import (
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo"
)

type AccountUsecase interface {
	CreateAccount(acc *domain.Account) error
	GetAllAccounts() ([]*domain.Account, error)
	GetAccountByID(accountID uint64) (*domain.Account, error)
	DeleteAccount(accountID uint64) error
	UpdateAccount(acc *domain.Account) error
}

type IntegrationUsecase interface {
	CreateIntegration(in *domain.Integration) error
	GetIntegrationsByAccountID(accountID uint64) ([]*domain.Integration, error)
}

type accountUsecase struct {
	repo repo.Repository
}

type integrationUsecase struct {
	repo repo.Repository
}

func NewAccountUsecase(r repo.Repository) AccountUsecase {
	return &accountUsecase{repo: r}
}

func NewIntegrationUsecase(r repo.Repository) IntegrationUsecase {
	return &integrationUsecase{repo: r}
}

type ContactsUsecase interface {
	SaveContacts(accountID uint64, contacts []*domain.Contact) error
	GetContactsByAccountID(accountID uint64) ([]*domain.Contact, error)
}

type contactsUsecase struct {
	repo repo.Repository
}

func NewContactsUsecase(r repo.Repository) ContactsUsecase {
	return &contactsUsecase{repo: r}
}
