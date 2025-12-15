package api

import (
	"encoding/json"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)

func (api *apiConfig) HandleCreateIntegration(w http.ResponseWriter, r *http.Request) {
	var in domain.Integration

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if in.AccountID == 0 {
		respondWithError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	if err := api.integrationUC.CreateIntegration(&in); err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot create integration")
		return
	}

	respondWithJSON(w, http.StatusCreated, in)
}
