package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "URL Path: %s", r.URL.Path[1:])
		})

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("Get failed with error: %s", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body with error: %s", err)
	}
	defer resp.Body.Close()

	bodyString := string(body)
	if bodyString != "URL Path: test" {
		t.Errorf("expected 'URL Path: test', got: %s", bodyString)
	}

}
