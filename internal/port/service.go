package port

import (
	"errors"
)

var ErrNotFound = errors.New("Port not found")

type Port struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

type repository interface {
	Exists(portID string) error
	Create(portID string, port Port) error
	Update(portID string, port Port) error
}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Upsert(portID string, port Port) error {
	if err := s.repo.Exists(portID); err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}
		return s.repo.Create(portID, port)
	}

	return s.repo.Update(portID, port)
}
