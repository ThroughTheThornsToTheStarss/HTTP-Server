package api

import (
	"log"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
)


func (api *apiConfig) HandleAmoAuthStart(w http.ResponseWriter, r *http.Request) {
	state := "amo_oauth_state"

	authURL := api.amoClient.AuthURL(state)
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

	tokens, err := api.amoClient.ExchangeCode(ctx, code)
	if err != nil {
		log.Println("exchange error:", err)
		respondWithError(w, http.StatusInternalServerError, "token exchange failed")
		return
	}

	acc := &domain.Account{
		ID:           api.amoClient.BaseDomain(),
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
