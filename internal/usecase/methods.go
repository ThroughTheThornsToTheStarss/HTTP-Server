package usecase

import (
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)

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

func (u *integrationUsecase) CreateIntegration(in *domain.Integration) error {
	return u.repo.CreateIntegration(in)
}

func (u *integrationUsecase) GetIntegrationsByAccountID(accountID string) ([]*domain.Integration, error) {
	return u.repo.GetIntegrationsByAccountID(accountID)
}
