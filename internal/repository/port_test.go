package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/rdleal/ports/internal/port"
)

func TestPort_Exists(t *testing.T) {
	testCases := []struct {
		name    string
		db      map[string]port.Port
		portID  string
		wantErr error
	}{
		{
			name:    "PortNotFound",
			db:      make(map[string]port.Port),
			portID:  "some non-existent port ID",
			wantErr: port.ErrNotFound,
		},
		{
			name:    "PortFound",
			db:      map[string]port.Port{"AEAJM": port.Port{Name: "Ajman"}},
			portID:  "AEAJM",
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewPort(tc.db)

			if got := repo.Exists(tc.portID); !errors.Is(got, tc.wantErr) {
				t.Errorf("got error calling Exists(%q): %v; want %v", tc.portID, got, tc.wantErr)
			}
		})
	}
}

func TestPort_Create(t *testing.T) {
	db := make(map[string]port.Port)
	repo := NewPort(db)

	portID := "AEAUH"
	port := port.Port{Name: "Abu Dhabi"}

	if err := repo.Create(portID, port); err != nil {
		t.Fatalf("got unexpected error calling Create(%q, %v): %q", portID, port, err)
	}

	gotPort, ok := db[portID]
	if !ok {
		t.Fatalf("got port with ID %q not found in db", portID)
	}

	if wantPort := port; !reflect.DeepEqual(gotPort, wantPort) {
		t.Errorf("got created port with ID %q: %v; want %v", portID, gotPort, wantPort)
	}
}

func TestPort_Update(t *testing.T) {
	db := map[string]port.Port{
		"AEAJM": port.Port{Name: "Dubai"},
	}
	repo := NewPort(db)

	portID := "AEAJM"
	port := port.Port{Name: "Ajman"}

	if err := repo.Update(portID, port); err != nil {
		t.Fatalf("got unexpected error calling Update(%q, %v): %q", portID, port, err)
	}

	gotPort, ok := db[portID]
	if !ok {
		t.Fatalf("got port with ID %q not found in db", portID)
	}

	if wantPort := port; !reflect.DeepEqual(gotPort, wantPort) {
		t.Errorf("got port with ID %q:  %v; want %v", portID, gotPort, wantPort)
	}
}
