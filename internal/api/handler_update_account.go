package api

import (
	"encoding/json"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)

func (api *apiConfig) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	var acc domain.Account

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if acc.ID == "" {
		respondWithError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	if err := api.accountUC.UpdateAccount(&acc); err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot update account")
		return
	}

	respondWithJSON(w, http.StatusOK, acc)
}
