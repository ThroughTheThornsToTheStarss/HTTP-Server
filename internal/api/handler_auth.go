package api

import (
	"log"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)


func (api *apiConfig) HandleAmoAuthStart(w http.ResponseWriter, r *http.Request) {
	authURL := api.amoClient.AuthURL()
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (api *apiConfig) HandleAmoAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if errStr := r.URL.Query().Get("error"); errStr != "" {
		respondWithError(w, http.StatusBadRequest, "amo error: "+errStr)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		respondWithError(w, http.StatusBadRequest, "code is missing")
		return
	}

	referer := r.URL.Query().Get("referer")
	if referer == "" {
		respondWithError(w, http.StatusBadRequest, "referer is missing")
		return
	}

	tokens, err := api.amoClient.ExchangeCode(ctx, referer, code)
	if err != nil {
		log.Println("exchange error:", err)
		respondWithError(w, http.StatusInternalServerError, "token exchange failed")
		return
	}

	acc := &domain.Account{
		ID:           referer,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Expires:      tokens.ExpiresAt,
	}

	if err := api.accountUC.CreateAccount(acc); err != nil {
		log.Println("save account error:", err)
		respondWithError(w, http.StatusInternalServerError, "cannot save authorization")
		return
	}

	respondWithJSON(w, http.StatusOK, acc)
}
