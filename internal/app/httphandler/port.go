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

	file, _, err := req.FormFile("ports")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	// checks the first object delimiter
	_, err = dec.Token()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	for dec.More() {
		portID, err := dec.Token()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		var p port.Port
		if err = dec.Decode(&p); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := h.service.Upsert(portID.(string), p); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}
