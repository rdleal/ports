package httphandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
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

func assertResponseStatus(t *testing.T, resp *http.Response, wantStatus int) {
	t.Helper()

	if gotStatus := resp.StatusCode; gotStatus != wantStatus {
		t.Errorf("got status code in response: %q; want %q",
			http.StatusText(gotStatus), http.StatusText(wantStatus))
	}

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

		assertResponseStatus(t, w.Result(), http.StatusOK)

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

func TestPort_ServeHTTP_Error(t *testing.T) {
	t.Run("InvalidMultipart", func(t *testing.T) {
		req, err := http.NewRequest("POST", "port-service-url", &bytes.Buffer{})
		if err != nil {
			t.Fatalf("got unexpected error creating http request: %v", err)
		}

		w := httptest.NewRecorder()

		service := stubPortServiceFunc(func(portID string, p port.Port) error { return nil })

		handler := NewPort(service)

		handler.ServeHTTP(w, req)

		assertResponseStatus(t, w.Result(), http.StatusBadRequest)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJSONFmt :=
			`--xxx 
Content-Disposition: form-data; name="ports"; filename="ports.json"
Content-Type: application/json 

%s
--xxx--
`

		testCases := []struct {
			name    string
			content string
		}{
			{
				name:    "ObjectDelimiter",
				content: "invalid-json",
			},
			{
				name:    "PortID",
				content: "{invalid-token",
			},
			{
				name:    "Port",
				content: `{"PORT_ID": {invalid-port-token}}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				payload := fmt.Sprintf(invalidJSONFmt, tc.content)

				body := io.NopCloser(strings.NewReader(payload))
				req, err := http.NewRequest("POST", "port-service-url", body)
				if err != nil {
					t.Fatalf("got unexpected error creating http request: %v", err)
				}

				req.Header.Set("Content-Type", "multipart/form-data; boundary=xxx")

				w := httptest.NewRecorder()

				service := stubPortServiceFunc(func(portID string, p port.Port) error { return nil })

				handler := NewPort(service)

				handler.ServeHTTP(w, req)

				assertResponseStatus(t, w.Result(), http.StatusBadRequest)

			})
		}

	})

	t.Run("ServiceError", func(t *testing.T) {
		path := "testdata/ports.json"

		body, multi := multipartBodyFromFile(t, path)

		req, err := http.NewRequest("POST", "port-service-url", body)
		if err != nil {
			t.Fatalf("got unexpected error creating http request: %v", err)
		}

		req.Header.Set("Content-Type", multi.FormDataContentType())

		w := httptest.NewRecorder()

		service := stubPortServiceFunc(func(portID string, p port.Port) error {
			return errors.New("some service error")
		})

		handler := NewPort(service)

		handler.ServeHTTP(w, req)

		assertResponseStatus(t, w.Result(), http.StatusInternalServerError)
	})
}
