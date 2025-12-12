package domain

type Contact struct {
	ID        uint    `json:"id"`
	AccountID string  `json:"account_id"`
	Name      string  `json:"name"`
	Email     *string `json:"email,omitempty"`
}
