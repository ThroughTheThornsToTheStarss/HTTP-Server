package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	a, err := app.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		log.Printf("HTTP server on :%s", a.HTTPPort)
		if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http error: %v", err)
		}
	}()

	go func() {
		log.Printf("gRPC server on :%s", a.GRPCPort)
		if err := a.GRPCHandler.Run(); err != nil {
			log.Printf("grpc error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.HTTPServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}
	a.GRPCHandler.Stop()

	if a.Producer != nil {
		if err := a.Producer.Close(); err != nil {
			log.Printf("beanstalk producer close error: %v", err)
		}
	}

}
