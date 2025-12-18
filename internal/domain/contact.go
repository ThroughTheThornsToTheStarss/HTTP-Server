package domain

type Contact struct {
	ID        uint    `json:"id"`
	AccountID uint64  `json:"account_id"`
	Name      string  `json:"name"`
	Email     *string `json:"email,omitempty"`
	Status    string  `json:"status"`
	AmoID  int64  `json:"amo_id"`
}
