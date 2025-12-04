package in_memory

import (
	"fmt"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)


func (r *MemoryRepository) CreateAccount(acc *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts[acc.ID] = acc
	return nil
}

func (r *MemoryRepository) GetAllAccounts() ([]*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]*domain.Account, 0, len(r.accounts))
	for _, acc := range r.accounts {
		res = append(res, acc)
	}
	return res, nil
}

func (r *MemoryRepository) CreateIntegration(in *domain.Integration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.integrations[in.AccountID] = append(r.integrations[in.AccountID], in)
	return nil
}

func (r *MemoryRepository) GetIntegrationsByAccountID(accountID string) ([]*domain.Integration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list, ok := r.integrations[accountID]
	if !ok {
		return []*domain.Integration{}, nil
	}

	return list, nil
}

func (r *MemoryRepository) DeleteAccount(accountID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.accounts, accountID)
	delete(r.integrations, accountID)
	return nil
}

func (r *MemoryRepository) UpdateAccount(acc *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.accounts[acc.ID]; !ok {
		return fmt.Errorf("account with id %s not found", acc.ID)
	}

	r.accounts[acc.ID] = acc
	return nil
}
