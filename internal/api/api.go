package api

import "git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"

type apiConfig struct {
	accountUC     usecase.AccountUsecase
	integrationUC usecase.IntegrationUsecase
}

func NewAPI(accountUC usecase.AccountUsecase, integrationUC usecase.IntegrationUsecase) *apiConfig {
	return &apiConfig{
		accountUC:     accountUC,
		integrationUC: integrationUC,
	}
}
