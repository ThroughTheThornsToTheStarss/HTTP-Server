package api

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
)

var (
	reEventID   = regexp.MustCompile(`^contacts\[(add|update|restore)\]\[\d+\]\[id\]$`)
	reEventName = regexp.MustCompile(`^contacts\[(add|update|restore)\]\[\d+\]\[name\]$`)
	reIndex     = regexp.MustCompile(`\[(\d+)\]`)
	reDel1      = regexp.MustCompile(`^contacts\[delete\]$`)
	reDel2      = regexp.MustCompile(`^contacts\[delete\]\[\d+\]$`)
	reDel3      = regexp.MustCompile(`^contacts\[delete\]\[\d+\]\[id\]$`)
)

func (api *apiConfig) HandleContactsWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := r.ParseForm(); err != nil {
		respondWithError(w, http.StatusBadRequest, "cannot parse form")
		return
	}

	if api.webhookUC == nil {
		respondWithError(w, http.StatusInternalServerError, "webhook usecase is not configured")
		return
	}

	accountID, err := parseAccountID(r.PostForm)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	events := parseEvents(r.PostForm)
	if len(events) == 0 {
		respondWithJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}

	if api.producer == nil {
		respondWithError(w, http.StatusServiceUnavailable, "queue is not available")
		return
	}

	jobID, err := api.webhookUC.Handle(r.Context(), accountID, events)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "cannot handle webhook")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"job_id": jobID,
	})
}

func parseAccountID(form url.Values) (uint64, error) {
	for k, v := range form {
		if len(v) == 0 {
			continue
		}
		if strings.HasPrefix(k, "contacts[") && strings.HasSuffix(k, "[account_id]") {
			id, err := strconv.ParseUint(strings.TrimSpace(v[0]), 10, 64)
			if err == nil && id > 0 {
				return id, nil
			}
		}
	}

	if s := strings.TrimSpace(form.Get("account_id")); s != "" {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil || id == 0 {
			return 0, errors.New("invalid account_id")
		}
		return id, nil
	}

	return 0, errors.New("account_id not found")
}

func parseEvents(form url.Values) []usecase.WebhookContactEvent {
	type tmp struct {
		id   int64
		name string
	}

	items := map[int]tmp{}

	for k, v := range form {
		if len(v) == 0 {
			continue
		}
		val := strings.TrimSpace(v[0])

		if reEventID.MatchString(k) {
			nums := reIndex.FindAllStringSubmatch(k, -1)
			if len(nums) == 0 {
				continue
			}
			idx, _ := strconv.Atoi(nums[len(nums)-1][1])
			id, _ := strconv.ParseInt(val, 10, 64)

			t := items[idx]
			t.id = id
			items[idx] = t
			continue
		}

		if reEventName.MatchString(k) {
			nums := reIndex.FindAllStringSubmatch(k, -1)
			if len(nums) == 0 {
				continue
			}
			idx, _ := strconv.Atoi(nums[len(nums)-1][1])

			t := items[idx]
			t.name = val
			items[idx] = t
			continue
		}
	}

	out := make([]usecase.WebhookContactEvent, 0, len(items))

	idxs := make([]int, 0, len(items))
	for idx := range items {
		idxs = append(idxs, idx)
	}
	sort.Ints(idxs)

	for _, idx := range idxs {
		t := items[idx]
		if t.id <= 0 {
			continue
		}
		out = append(out, usecase.WebhookContactEvent{
			AmoID:   t.id,
			Name:    t.name,
			Email:   nil,
			Deleted: false,
		})
	}

	for _, id := range parseDeletedIDs(form) {
		out = append(out, usecase.WebhookContactEvent{
			AmoID:   id,
			Deleted: true,
		})
	}

	return out
}

func parseDeletedIDs(form url.Values) []int64 {
	set := map[int64]struct{}{}

	for k, v := range form {
		if len(v) == 0 {
			continue
		}
		val := strings.TrimSpace(v[0])

		if reDel1.MatchString(k) || reDel2.MatchString(k) || reDel3.MatchString(k) {
			id, err := strconv.ParseInt(val, 10, 64)
			if err == nil && id > 0 {
				set[id] = struct{}{}
			}
		}
	}

	if len(set) == 0 {
		return nil
	}

	out := make([]int64, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
