package in_memory

import (
	"sync"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)

type MemoryRepository struct {
	mu           sync.RWMutex
	accounts     map[string]*domain.Account
	integrations map[string][]*domain.Integration
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		accounts:     make(map[string]*domain.Account),
		integrations: make(map[string][]*domain.Integration),
	}
}
