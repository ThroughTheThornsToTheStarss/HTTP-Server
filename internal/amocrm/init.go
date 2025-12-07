package amocrm

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type tokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresAt    int64
}

type OAuthClient struct {
	httpClient   *http.Client
	baseDomain   string
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewOAuthClientFromEnv() (*OAuthClient, error) {
	baseDomain := os.Getenv("AMO_BASE_DOMAIN")
	clientID := os.Getenv("AMO_CLIENT_ID")
	clientSecret := os.Getenv("AMO_CLIENT_SECRET")
	redirectURI := os.Getenv("AMO_REDIRECT_URI")

	if baseDomain == "" || clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, fmt.Errorf("env not set")
	}

	return &OAuthClient{
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		baseDomain:   baseDomain,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}, nil
}
