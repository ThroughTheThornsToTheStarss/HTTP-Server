package app

import (
	"fmt"
	"net/http"
	"os"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/amocrm"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/api"
	grpcdelivery "git.amocrm.ru/ilnasertdinov/http-server-go/internal/grpc"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
	beanstalkq "git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue/beanstalk"
	mysqlrepo "git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo/mysql"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/usecase"
	mysqlcfg "git.amocrm.ru/ilnasertdinov/http-server-go/pkg/mysql"
)

type App struct {
	HTTPServer  *http.Server
	GRPCHandler *grpcdelivery.Handler
	HTTPPort    string
	GRPCPort    string
	Producer    queue.Producer
}

func NewFromEnv() (*App, error) {
	httpPort := getenv("HTTP_PORT", "8080")
	grpcPort := getenv("GRPC_PORT", "8091")

	db, err := mysqlcfg.NewGormFromEnv()
	if err != nil {
		return nil, fmt.Errorf("mysql connect: %w", err)
	}

	if err := mysqlrepo.AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("mysql migrate: %w", err)
	}

	repo := mysqlrepo.NewGormRepository(db)
	accountUC := usecase.NewAccountUsecase(repo)
	integrationUC := usecase.NewIntegrationUsecase(repo)
	contactsUC := usecase.NewContactsUsecase(repo)

	amoClient, err := amocrm.NewOAuthClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("amo oauth client: %w", err)
	}

	producer, err := beanstalkq.New(getenv("BEANSTALK_ADDR", "beanstalkd:11300"))
	if err != nil {
		return nil, fmt.Errorf("beanstalk connect: %w", err)
	}
	handler := api.New(accountUC, integrationUC, contactsUC, amoClient, producer)
	httpSrv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: handler,
	}

	grpcHandler, err := grpcdelivery.NewHandler(":"+grpcPort, grpcdelivery.NewAccountServer(accountUC))
	if err != nil {
		return nil, fmt.Errorf("grpc init: %w", err)
	}

	return &App{
		HTTPServer:  httpSrv,
		GRPCHandler: grpcHandler,
		HTTPPort:    httpPort,
		GRPCPort:    grpcPort,
		Producer:    producer,
	}, nil
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
