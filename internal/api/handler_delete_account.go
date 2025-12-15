package api

import (
	"net/http"
	"strconv"
)

func (api *apiConfig) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	accountIDInt, err := strconv.ParseUint(accountID, 10, 64)
	if err != nil || accountIDInt == 0 {
		respondWithError(w, http.StatusBadRequest, "account_id must be a positive integer")
		return
	}

	err = api.accountUC.DeleteAccount(accountIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot delete account")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
