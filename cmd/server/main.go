package main

import (
	"log"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/api"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo/in_memory"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
)

func main() {

	const port = "8080"

	repo := in_memory.NewMemoryRepository()
	accountUC := usecase.NewAccountUsecase(repo)
	integrationUC := usecase.NewIntegrationUsecase(repo)

	handler := api.New(accountUC, integrationUC)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
