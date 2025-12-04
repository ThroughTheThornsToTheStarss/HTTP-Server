package api

import (
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo"
)

type apiConfig struct {
	repo repo.Repository
}

func NewAPI(repo repo.Repository) *apiConfig {
	return &apiConfig{
		repo: repo,
	}
}
