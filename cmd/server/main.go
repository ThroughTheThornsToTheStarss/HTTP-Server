package main

import (
	"log"
	"net/http"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/amocrm"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/api"
	mysqlrepo "git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo/mysql"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
	mysqlcfg "git.amocrm.ru/ilnasertdinov/http-server-go/pkg/mysql"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	const port = "8080"

	db, err := mysqlcfg.NewGormFromEnv()
	if err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}

	if err := mysqlrepo.AutoMigrate(db); err != nil {
		log.Fatalf("mysql automigrate error: %v", err)
	}

	repo := mysqlrepo.NewGormRepository(db)
	accountUC := usecase.NewAccountUsecase(repo)
	integrationUC := usecase.NewIntegrationUsecase(repo)
	contactsUC := usecase.NewContactsUsecase(repo)

	amoClient, err := amocrm.NewOAuthClientFromEnv()
	if err != nil {
		log.Fatalf("cannot init amo oauth client: %v", err)
	}

	handler := api.New(accountUC, integrationUC, contactsUC, amoClient)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
