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

func NewAccountUsecase(r repo.Repository) AccountUsecase {
	return &accountUsecase{repo: r}
}

func (u *accountUsecase) CreateAccount(acc *domain.Account) error {
	return u.repo.CreateAccount(acc)
}

func (u *accountUsecase) GetAllAccounts() ([]*domain.Account, error) {
	return u.repo.GetAllAccounts()
}

func (u *accountUsecase) DeleteAccount(accountID string) error {
	return u.repo.DeleteAccount(accountID)
}

func (u *accountUsecase) UpdateAccount(acc *domain.Account) error {
	return u.repo.UpdateAccount(acc)
}

type integrationUsecase struct {
	repo repo.Repository
}

func NewIntegrationUsecase(r repo.Repository) IntegrationUsecase {
	return &integrationUsecase{repo: r}
}

func (u *integrationUsecase) CreateIntegration(in *domain.Integration) error {
	return u.repo.CreateIntegration(in)
}

func (u *integrationUsecase) GetIntegrationsByAccountID(accountID string) ([]*domain.Integration, error) {
	return u.repo.GetIntegrationsByAccountID(accountID)
}
