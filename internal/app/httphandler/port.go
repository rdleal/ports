package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rdleal/ports/internal/port"
)

type portService interface {
	Upsert(portID string, p port.Port) error
}

type Port struct {
	service portService
}

func NewPort(s portService) *Port {
	return &Port{service: s}
}

func (h *Port) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(200 << 20)

	file, _, _ := req.FormFile("ports")
	defer file.Close()

	ports := make(map[string]port.Port)
	json.NewDecoder(file).Decode(&ports)

	for portID, port := range ports {
		h.service.Upsert(portID, port)
	}
}
