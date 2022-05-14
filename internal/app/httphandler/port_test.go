package httphandler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/rdleal/ports/internal/port"
)

type stubPortServiceFunc func(portID string, p port.Port) error

func (f stubPortServiceFunc) Upsert(portID string, p port.Port) error {
	return f(portID, p)
}

func openTestFile(t *testing.T, path string) *os.File {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("got error openning the file %q: %v", path, err)
	}

	return f
}

func multipartBodyFromFile(t *testing.T, path string) (io.Reader, *multipart.Writer) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	f := openTestFile(t, path)

	part, err := writer.CreateFormFile("ports", path)
	if err != nil {
		t.Fatalf("got unexpected error creating multipart file %q: %v", path, err)
	}
	io.Copy(part, f)

	writer.Close()

	return body, writer
}

func TestPort_ServeHTTP(t *testing.T) {
	t.Run("UploadPorts", func(t *testing.T) {
		path := "testdata/ports.json"

		body, multi := multipartBodyFromFile(t, path)

		req, err := http.NewRequest("POST", "port-service-url", body)
		if err != nil {
			t.Fatalf("got unexpected error creating http request: %v", err)
		}

		req.Header.Set("Content-Type", multi.FormDataContentType())

		w := httptest.NewRecorder()

		gotPorts := make(map[string]port.Port)

		service := stubPortServiceFunc(func(portID string, p port.Port) error {
			gotPorts[portID] = p
			return nil
		})

		handler := NewPort(service)

		handler.ServeHTTP(w, req)

		if gotStatus, wantStatus := w.Result().StatusCode, http.StatusOK; gotStatus != wantStatus {
			t.Errorf("got status code in response: %q; want %q",
				http.StatusText(gotStatus), http.StatusText(wantStatus))
		}

		portsJSON := openTestFile(t, path)

		wantPorts := make(map[string]port.Port)

		if err := json.NewDecoder(portsJSON).Decode(&wantPorts); err != nil {
			t.Fatalf("got unexpected error parsing json file %q: %v", path, err)
		}

		if !reflect.DeepEqual(gotPorts, wantPorts) {
			t.Errorf("got ports in service.Upsert(): %v; want %v", gotPorts, wantPorts)
		}
	})
}
