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

	dec := json.NewDecoder(file)

	// checks the first object delimiter
	_, _ = dec.Token()

	for dec.More() {
		portID, _ := dec.Token()

		var p port.Port
		_ = dec.Decode(&p)

		h.service.Upsert(portID.(string), p)
	}
}
