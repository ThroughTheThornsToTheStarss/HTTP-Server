package main

import (
	"log"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/amocrm"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/api"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo/in_memory"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
	"git.amocrm.ru/ilnasertdinov/http-server-go/pkg/mysql"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	const port = "8080"


	db, err := mysql.ConnectFromEnv()
	if err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}
	defer db.Close()



	repo := in_memory.NewMemoryRepository()
	accountUC := usecase.NewAccountUsecase(repo)
	integrationUC := usecase.NewIntegrationUsecase(repo)


	amoClient, err := amocrm.NewOAuthClientFromEnv()
	if err != nil {
		log.Fatalf("cannot init amo oauth client: %v", err)
	}

	handler := api.New(accountUC, integrationUC, amoClient)

	srv := &http.Server{
		Addr:    ":" + port,

		Handler: handler,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
