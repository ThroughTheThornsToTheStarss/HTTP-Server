package amocrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func normalizeBaseDomain(baseDomain string) string {
	baseDomain = strings.TrimSpace(baseDomain)
	baseDomain = strings.TrimRight(baseDomain, "/")

	if !strings.HasPrefix(baseDomain, "http://") && !strings.HasPrefix(baseDomain, "https://") {
		baseDomain = "https://" + baseDomain
	}

	return baseDomain
}

func (c *OAuthClient) AuthURL() string {

	values := url.Values{}
	values.Set("client_id", c.clientID)
	values.Set("mode", "post_message")

	return fmt.Sprintf("%s/oauth?%s", c.baseDomain, values.Encode())
}

func (c *OAuthClient) ExchangeCode(ctx context.Context, baseDomain, code string) (Tokens, error) {
	baseDomain = normalizeBaseDomain(baseDomain)

	body := map[string]any{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  c.redirectURI,
	}

	return c.sendTokenRequest(ctx, baseDomain, body)
}

func (c *OAuthClient) sendTokenRequest(ctx context.Context, baseDomain string, body map[string]any) (Tokens, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return Tokens{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseDomain+"/oauth2/access_token",
		bytes.NewReader(data),
	)
	if err != nil {
		return Tokens{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Tokens{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return Tokens{}, fmt.Errorf("amo error: status=%d body=%s", resp.StatusCode, errBody)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		TokenType:    tr.TokenType,
		ExpiresIn:    tr.ExpiresIn,    
		ExpiresAt:    time.Now().Unix() + tr.ExpiresIn,
	}, nil
}
