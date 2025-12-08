package api

import (
	"log"
	"net/http"
)

type contactResponse struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
}

func (api *apiConfig) HandleAmoGetContacts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	accounts, err := api.accountUC.GetAllAccounts()
	if err != nil {
		log.Println("get accounts error:", err)
		respondWithError(w, http.StatusInternalServerError, "cannot load accounts")
		return
	}

	if len(accounts) == 0 {
		respondWithError(w, http.StatusBadRequest, "no authorized amo accounts")
		return
	}

	acc := accounts[0]

	amoContacts, err := api.amoClient.GetAllContacts(ctx, acc.ID, acc.AccessToken)
	if err != nil {
		log.Println("amo get all contacts error:", err)
		respondWithError(w, http.StatusInternalServerError, "cannot fetch contacts from amo")
		return
	}

	resp := make([]contactResponse, 0, len(amoContacts))

	for _, c := range amoContacts {
		item := contactResponse{
			Name: c.Name,
		}

		if email, ok := c.PrimaryEmail(); ok {
			item.Email = &email
		} else {
			item.Email = nil
		}

		resp = append(resp, item)
	}

	respondWithJSON(w, http.StatusOK, resp)
}
