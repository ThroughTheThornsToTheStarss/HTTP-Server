package api

import (
	"log"
	"net/http"
	"strconv"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)


func (api *apiConfig) HandleAmoGetContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountIDStr := r.URL.Query().Get("account_id")

	if accountIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid account_id")
		return
	}

	acc, err := api.accountUC.GetAccountByID(accountID)
	if err != nil {
		log.Println("get account error:", err)
		respondWithError(w, http.StatusNotFound, "account not found")
		return
	}

	amoContacts, err := api.amoClient.GetAllContacts(ctx, acc.Referer, acc.AccessToken)
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
