package repository

import (
	"github.com/rdleal/ports/internal/port"
)

type Port struct {
	db map[string]port.Port
}

func NewPort(db map[string]port.Port) *Port {
	return &Port{db: db}
}

func (r *Port) Exists(portID string) error {
	if _, ok := r.db[portID]; !ok {
		return port.ErrNotFound
	}

	return nil
}

func (r *Port) Create(portID string, p port.Port) error {
	r.db[portID] = p
	return nil
}

func (r *Port) Update(portID string, p port.Port) error {
	r.db[portID] = p
	return nil
}
