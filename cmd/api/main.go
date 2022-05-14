package main

import (
	"fmt"

	"github.com/rdleal/ports/internal/port"
	"github.com/rdleal/ports/internal/repository"
)

func main() {

	db := make(map[string]port.Port)
	repo := repository.NewPort(db)
	service := port.NewService(repo)

	service.Upsert("SomePort", port.Port{Name: "Some Name"})

	fmt.Println(db)
}
