package api

import (
	"net/http"
)

func (api *apiConfig) HandleGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := api.accountUC.GetAllAccounts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot get accounts")
		return
	}

	respondWithJSON(w, http.StatusOK, accounts)
}
