package amocrm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c Contact) PrimaryEmail() (string, bool) {
	for _, f := range c.CustomFields {
		if f.FieldCode == "EMAIL" && len(f.Values) > 0 {
			return f.Values[0].Value, true
		}
	}
	return "", false
}


func (c *OAuthClient) getContactsPage(ctx context.Context, baseDomain string, accessToken string, page, limit int) ([]Contact, error) {
	var result contactsPage

	baseDomain = normalizeBaseDomain(baseDomain)

	u, err := url.Parse(baseDomain + "/api/v4/contacts")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer " + accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status=%d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embedded.Contacts, nil
}

func (c *OAuthClient) GetAllContacts(ctx context.Context, baseDomain string, accessToken string) ([]Contact, error) {
	const pageLimit = 250

	var all []Contact
	page := 1

	for {
		contacts, err := c.getContactsPage(ctx, baseDomain, accessToken, page, pageLimit)
		if err != nil {
			return nil, err
		}

		if len(contacts) == 0 {
			break
		}
		all = append(all, contacts...)

		if len(contacts) < pageLimit {
			break
		}
		page++
	}

	return all, nil
}
