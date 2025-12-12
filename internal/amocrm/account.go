package amocrm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *OAuthClient) GetAccountInfo(ctx context.Context, baseDomain, accessToken string) (AccountInfo, error) {
	baseDomain = normalizeBaseDomain(baseDomain)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseDomain+"/api/v4/account", nil)
	if err != nil {
		return AccountInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AccountInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AccountInfo{}, fmt.Errorf("amo account info error: status=%d", resp.StatusCode)
	}

	var info AccountInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return AccountInfo{}, err
	}
	if info.ID == 0 {
		return AccountInfo{}, fmt.Errorf("amo account info: empty id")
	}

	return info, nil
}
