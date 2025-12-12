package domain

type Integration struct {
	AccountID          uint64 `json:"account_id"`
	SecretKey          string `json:"secret_key"`
	ClientID           string `json:"client_id"`
	RedirectURL        string `json:"redirect_url"`
	AuthenticationCode string `json:"authentication_code"`
}
