package main

import (
	"log"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/api"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repository/in_memory"
)

func main() {

	const port = "8080"

	repo := in_memory.NewMemoryRepository()
	apiCfg := api.NewAPI(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /accounts", apiCfg.HandlerCreateAccount)
	mux.HandleFunc("GET /accounts", apiCfg.HandleGetAllAccounts)

	mux.HandleFunc("POST /integrations", apiCfg.HandleCreateIntegration)
	mux.HandleFunc("GET /integrations", apiCfg.HandleGetIntegrations)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
