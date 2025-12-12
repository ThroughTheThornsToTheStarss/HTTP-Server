package api

import (
	"log"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
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

	contacts := make([]*domain.Contact, 0, len(amoContacts))
	for _, c := range amoContacts {
		var emailPtr *string
		if email, ok := c.PrimaryEmail(); ok {
			emailCopy := email
			emailPtr = &emailCopy
		}

		contacts = append(contacts, &domain.Contact{
			AccountID: acc.ID,
			Name:      c.Name,
			Email:     emailPtr,
		})
	}

	if err := api.contactsUC.SaveContacts(acc.ID, contacts); err != nil {
		log.Println("save contacts error:", err)
		respondWithError(w, http.StatusInternalServerError, "cannot save contacts")
		return
	}

	resp := make([]contactResponse, 0, len(contacts))
	for _, c := range contacts {
		resp = append(resp, contactResponse{
			Name:  c.Name,
			Email: c.Email,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}
