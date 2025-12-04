package api

import (
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
)

type apiConfig struct {
	accountUC     usecase.AccountUsecase
	integrationUC usecase.IntegrationUsecase
}

func New(accountUC usecase.AccountUsecase, 	integrationUC usecase.IntegrationUsecase) http.Handler {
	apiCfg := &apiConfig{
		accountUC:     accountUC,
		integrationUC: integrationUC,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /accounts", apiCfg.HandlerCreateAccount)
	mux.HandleFunc("GET /accounts", apiCfg.HandleGetAllAccounts)
	mux.HandleFunc("DELETE /accounts", apiCfg.HandleDeleteAccount)
	mux.HandleFunc("PUT /accounts", apiCfg.HandleUpdateAccount)

	mux.HandleFunc("POST /integrations", apiCfg.HandleCreateIntegration)
	mux.HandleFunc("GET /integrations", apiCfg.HandleGetIntegrations)

	return mux
}
