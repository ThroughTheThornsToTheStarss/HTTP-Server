package usecase

import (
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo"
)

type AccountUsecase interface {
	CreateAccount(acc *domain.Account) error
	GetAllAccounts() ([]*domain.Account, error)
	DeleteAccount(accountID string) error
	UpdateAccount(acc *domain.Account) error
}

type IntegrationUsecase interface {
	CreateIntegration(in *domain.Integration) error
	GetIntegrationsByAccountID(accountID string) ([]*domain.Integration, error)
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
	SaveContacts(accountID string, contacts []*domain.Contact) error
	GetContactsByAccountID(accountID string) ([]*domain.Contact, error)
}

type contactsUsecase struct {
	repo repo.Repository
}

func NewContactsUsecase(r repo.Repository) ContactsUsecase {
	return &contactsUsecase{repo: r}
}
