package domain

type Repository interface {

	CreateAccount(acc *Account) error
	GetAllAccounts() ([]*Account, error)

	CreateIntegration(in *Integration) error
	GetIntegrationsByAccountID(accountID string) ([]*Integration, error)
}

