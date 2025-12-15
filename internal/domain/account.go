package domain

type Account struct {
	ID           uint64 `json:"id"`
	Referer      string `json:"referer"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
