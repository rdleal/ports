package main

import (
	"fmt"

	"github.com/rdleal/ports/internal/port"
)

type fakeRepo struct {
}

func (r *fakeRepo) Exists(portID string) error {
	fmt.Printf("Exists(%q)\n", portID)
	return nil
}

func (r *fakeRepo) Create(portID string, p port.Port) error {
	fmt.Printf("Create(%q, %v)\n", portID, p)
	return nil
}

func (r *fakeRepo) Update(portID string, p port.Port) error {
	fmt.Printf("Update(%q, %v)\n", portID, p)
	return nil
}

func main() {

	repo := &fakeRepo{}
	service := port.NewService(repo)

	fmt.Println(service.Upsert("SomePort", port.Port{Name: "Some Name"}))
}
