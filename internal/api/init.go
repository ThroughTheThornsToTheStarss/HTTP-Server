package api

import (
	"net/http"


	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/amocrm"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
)
type contactResponse struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
}

type apiConfig struct {
	accountUC     usecase.AccountUsecase
	integrationUC usecase.IntegrationUsecase
	contactsUC    usecase.ContactsUsecase
	amoClient     *amocrm.OAuthClient
}

func New(accountUC usecase.AccountUsecase, integrationUC usecase.IntegrationUsecase, contactsUC usecase.ContactsUsecase, amoClient *amocrm.OAuthClient) http.Handler {
	apiCfg := &apiConfig{
		accountUC:     accountUC,
		integrationUC: integrationUC,
		contactsUC:    contactsUC,
		amoClient:     amoClient,
=======
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
)



	mux := http.NewServeMux()

	mux.HandleFunc("POST /accounts", apiCfg.HandlerCreateAccount)
	mux.HandleFunc("GET /accounts", apiCfg.HandleGetAllAccounts)

	mux.HandleFunc("DELETE /accounts", apiCfg.HandleDeleteAccount)
=======


	mux.HandleFunc("PUT /accounts", apiCfg.HandleUpdateAccount)

	mux.HandleFunc("POST /integrations", apiCfg.HandleCreateIntegration)
	mux.HandleFunc("GET /integrations", apiCfg.HandleGetIntegrations)


	mux.HandleFunc("GET /amo/auth/start", apiCfg.HandleAmoAuthStart)
	mux.HandleFunc("GET /amo/oauth/callback", apiCfg.HandleAmoAuthCallback)

	mux.HandleFunc("GET /amo/contacts", apiCfg.HandleAmoGetContacts)


	return mux
}
