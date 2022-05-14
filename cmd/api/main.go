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

	http.Handle("/ports", handler)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
