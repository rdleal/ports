package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rdleal/ports/internal/app/httphandler"
	"github.com/rdleal/ports/internal/port"
	"github.com/rdleal/ports/internal/repository"
)

func main() {
	db := make(map[string]port.Port)
	repo := repository.NewPort(db)
	service := port.NewService(repo)
	handler := httphandler.NewPort(service)

	httpPort := os.Getenv("PORT")

	log.Printf("Running service on port: %s\n", httpPort)

	http.Handle("/ports", handler)
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		log.Fatal(err)
	}
}
