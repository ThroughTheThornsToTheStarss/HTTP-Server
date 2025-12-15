package usecase

import "git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"

func (u *accountUsecase) CreateAccount(acc *domain.Account) error {
	return u.repo.CreateAccount(acc)
}

func (u *accountUsecase) GetAllAccounts() ([]*domain.Account, error) {
	return u.repo.GetAllAccounts()
}

func (u *accountUsecase) GetAccountByID(accountID uint64) (*domain.Account, error) {
	return u.repo.GetAccountByID(accountID)
}

func (u *accountUsecase) DeleteAccount(accountID uint64) error {
	return u.repo.DeleteAccount(accountID)
}

func (u *accountUsecase) UpdateAccount(acc *domain.Account) error {
	return u.repo.UpdateAccount(acc)
}

func (u *integrationUsecase) CreateIntegration(in *domain.Integration) error {
	return u.repo.CreateIntegration(in)
}

func (u *integrationUsecase) GetIntegrationsByAccountID(accountID uint64) ([]*domain.Integration, error) {
	return u.repo.GetIntegrationsByAccountID(accountID)
}

func (u *contactsUsecase) SaveContacts(accountID uint64, contacts []*domain.Contact) error {
	return u.repo.SaveContacts(accountID, contacts)
}

func (u *contactsUsecase) GetContactsByAccountID(accountID uint64) ([]*domain.Contact, error) {
	return u.repo.GetContactsByAccountID(accountID)
}
