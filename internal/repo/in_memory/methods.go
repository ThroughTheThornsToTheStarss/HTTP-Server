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

func (r *MemoryRepository) GetIntegrationsByAccountID(accountID uint64) ([]*domain.Integration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list, ok := r.integrations[accountID]
	if !ok {
		return []*domain.Integration{}, nil
	}

	return list, nil
}

func (r *MemoryRepository) DeleteAccount(accountID uint64) error {
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
		return fmt.Errorf("account with id %d not found", acc.ID)
	}

	r.accounts[acc.ID] = acc
	return nil
}

func (r *MemoryRepository) SaveContacts(accountID uint64, contacts []*domain.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := make([]*domain.Contact, len(contacts))
	copy(cp, contacts)
	r.contacts[accountID] = cp
	return nil
}

func (r *MemoryRepository) GetContactsByAccountID(accountID uint64) ([]*domain.Contact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	src := r.contacts[accountID]
	if src == nil {
		return []*domain.Contact{}, nil
	}

	cp := make([]*domain.Contact, len(src))
	copy(cp, src)
	return cp, nil
}
