package api

import (
	"net/http"
)

func (api *apiConfig) HandlrDeleteAccounts(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	err := api.repo.DeleteAccount(accountID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot delete account")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
