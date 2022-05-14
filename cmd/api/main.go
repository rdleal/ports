package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/rdleal/ports/internal/app/config"
	"github.com/rdleal/ports/internal/app/httphandler"
	"github.com/rdleal/ports/internal/port"
	"github.com/rdleal/ports/internal/repository"
)

func main() {
	db := make(map[string]port.Port)
	repo := repository.NewPort(db)
	service := port.NewService(repo)
	handler := httphandler.NewPort(service)

	httpPort := ":" + os.Getenv("PORT")

	ctx := config.ContextWithGracefulCancellation(context.Background())

	srv := &http.Server{
		Addr:    httpPort,
		Handler: handler,
	}

	log.Printf("Running service on port %s", httpPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe() error: %q", err)
		}
	}()

	<-ctx.Done()

	srv.Shutdown(context.Background())
}
