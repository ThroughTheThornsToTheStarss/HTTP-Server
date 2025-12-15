package in_memory

import (
	"sync"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)

type MemoryRepository struct {
	mu           sync.RWMutex
	accounts     map[uint64]*domain.Account
	integrations map[uint64][]*domain.Integration
	contacts     map[uint64][]*domain.Contact
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		accounts:     make(map[uint64]*domain.Account),
		integrations: make(map[uint64][]*domain.Integration),
		contacts:     make(map[uint64][]*domain.Contact),
	}
}
