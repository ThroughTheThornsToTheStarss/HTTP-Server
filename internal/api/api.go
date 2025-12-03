package api

import "git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"


type apiConfig struct {
	repo domain.Repository
}


func NewAPI(repo domain.Repository) *apiConfig {
	return &apiConfig{
		repo: repo,
	}
}
