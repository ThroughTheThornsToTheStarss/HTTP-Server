package api

import (
	"net/http"
)

func (api *apiConfig) HandleGetIntegrations(w http.ResponseWriter, r *http.Request) {

	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	list, err := api.repo.GetIntegrationsByAccountID(accountID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot get integrations")
		return
	}

	respondWithJSON(w, http.StatusOK, list)
}
