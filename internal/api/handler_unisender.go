package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)

func (api *apiConfig) HandleUnisenderKey(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if err := r.ParseForm(); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid form")
		return
	}

	key := strings.TrimSpace(r.FormValue("unisender_key"))
	if key == "" {
		key = strings.TrimSpace(r.FormValue("api_token"))
	}

	accStr := strings.TrimSpace(r.FormValue("account_id"))
	accID, err := strconv.ParseUint(accStr, 10, 64)
	if err != nil || accID == 0 {
		respondWithError(w, http.StatusBadRequest, "account_id must be a positive integer")
		return
	}
	if key == "" {
		respondWithError(w, http.StatusBadRequest, "unisender_key is required")
		return
	}

	in := &domain.Integration{
		AccountID:    accID,
		UnisenderKey: key,
	}
	if err := api.integrationUC.CreateIntegration(in); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]any{"ok": true})
}
