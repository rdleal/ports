package port

import (
	"errors"
	"reflect"
	"testing"
)

type stubRepo struct {
	existsReceivedPortID string
	existsReturnsErr     error

	createReceivedPortID string
	createReceivedPort   Port
	createReturnsErr     error

	updateReceivedPortID string
	updateReceivedPort   Port
	updateReturnsErr     error
}

func (r *stubRepo) Exists(portID string) error {
	r.existsReceivedPortID = portID

	return r.existsReturnsErr
}

func (r *stubRepo) Create(portID string, port Port) error {
	r.createReceivedPortID = portID
	r.createReceivedPort = port

	return r.createReturnsErr
}

func (r *stubRepo) Update(portID string, port Port) error {
	r.updateReceivedPortID = portID
	r.updateReceivedPort = port

	return r.updateReturnsErr
}

func TestService_Upsert(t *testing.T) {
	t.Run("CreatesPort", func(t *testing.T) {
		repo := &stubRepo{
			existsReturnsErr: ErrNotFound,
		}

		service := NewService(repo)

		portID := "AEAJM"
		port := Port{Name: "Ajman"}

		if err := service.Upsert(portID, port); err != nil {
			t.Fatalf("got unexpected error calling Upsert(%q, %v): %s", portID, port, err)
		}

		if got, want := repo.existsReceivedPortID, portID; got != want {
			t.Errorf("repo.Exists(): got port ID: %q; want %q", got, want)
		}

		if got, want := repo.createReceivedPortID, portID; got != want {
			t.Errorf("repo.Create(): got port ID: %q; want %q", got, want)
		}

		if got, want := repo.createReceivedPort, port; !reflect.DeepEqual(got, want) {
			t.Errorf("repo.Create(): got port: %v; want %v", got, want)
		}
	})

	t.Run("UpdatesPort", func(t *testing.T) {
		repo := &stubRepo{
			existsReturnsErr: nil,
		}

		service := NewService(repo)

		portID := "AEAJM"
		port := Port{Name: "Ajman"}

		if err := service.Upsert(portID, port); err != nil {
			t.Fatalf("got unexpected error calling Upsert(%q, %v): %s", portID, port, err)
		}

		if got, want := repo.existsReceivedPortID, portID; got != want {
			t.Errorf("repo.Exists(): got port ID: %q; want %q", got, want)
		}

		if got, want := repo.updateReceivedPortID, portID; got != want {
			t.Errorf("repo.Update(): got port ID: %q; want %q", got, want)
		}

		if got, want := repo.updateReceivedPort, port; !reflect.DeepEqual(got, want) {
			t.Errorf("repo.Update(): got port: %v; want %v", got, want)
		}
	})
}

func TestService_Upsert_Error(t *testing.T) {
	t.Run("ExistsError", func(t *testing.T) {
		repo := &stubRepo{
			existsReturnsErr: errors.New("repository exists error"),
		}

		service := NewService(repo)

		portID := "AEAJM"
		port := Port{Name: "Ajman"}

		got := service.Upsert(portID, port)
		if got == nil {
			t.Fatalf("got nil error calling Upsert(%q, %v); want not nil", portID, port)
		}

		if want := repo.existsReturnsErr; !errors.Is(got, want) {
			t.Errorf("got Upsert() error: %v; want %v", got, want)
		}
	})

	t.Run("CreateError", func(t *testing.T) {
		repo := &stubRepo{
			existsReturnsErr: ErrNotFound,
			createReturnsErr: errors.New("repository create error"),
		}

		service := NewService(repo)

		portID := "AEAJM"
		port := Port{Name: "Ajman"}

		got := service.Upsert(portID, port)
		if got == nil {
			t.Fatalf("got nil error calling Upsert(%q, %v); want not nil", portID, port)
		}

		if want := repo.createReturnsErr; !errors.Is(got, want) {
			t.Errorf("got Upsert() error: %v; want %v", got, want)
		}
	})

	t.Run("CreateError", func(t *testing.T) {
		repo := &stubRepo{
			existsReturnsErr: nil,
			updateReturnsErr: errors.New("repository update error"),
		}

		service := NewService(repo)

		portID := "AEAJM"
		port := Port{Name: "Ajman"}

		got := service.Upsert(portID, port)
		if got == nil {
			t.Fatalf("got nil error calling Upsert(%q, %v); want not nil", portID, port)
		}

		if want := repo.updateReturnsErr; !errors.Is(got, want) {
			t.Errorf("got Upsert() error: %v; want %v", got, want)
		}
	})
}
