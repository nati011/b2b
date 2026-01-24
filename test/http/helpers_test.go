package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const expectedStatusFmt = "expected status %d, got %d"

func doJSONRequest(t *testing.T, client *http.Client, method, url string, payload any) *http.Response {
	t.Helper()
	return doJSONRequestWithHeaders(t, client, method, url, payload, nil)
}

func doJSONRequestWithHeaders(t *testing.T, client *http.Client, method, url string, payload any, headers map[string]string) *http.Response {
	t.Helper()

	var bodyReader *bytes.Buffer
	if payload != nil {
		bodyReader = &bytes.Buffer{}
		if err := json.NewEncoder(bodyReader).Encode(payload); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}

	var body io.Reader
	if bodyReader != nil {
		body = bodyReader
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("perform request: %v", err)
	}
	return resp
}
