package in_memory

func (r *MemoryRepository) UpsertContactFromWebhook(accountID uint64, amoID int64, name string, email *string, status string) (uint, error) {
	return uint(amoID), nil
}

func (r *MemoryRepository) FindContactIDByAmoID(accountID uint64, amoID int64) (uint, bool, error) {
	return 0, false, nil
}

func (r *MemoryRepository) AddSyncHistory(contactID uint, status string, message string) error {
	return nil
}

func (r *MemoryRepository) TrimSyncHistory(contactID uint, keepLast int) error {
	return nil
}
