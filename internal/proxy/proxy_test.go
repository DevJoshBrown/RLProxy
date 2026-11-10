package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/DevJoshBrown/RLProxy/internal/algorithm"
)

func TestReverseProxy(t *testing.T) {

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "test message")
		})

	server := httptest.NewServer(handler)
	defer server.Close()

	rpURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse rpURL")
	}

	ReverseProxy := httptest.NewServer(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(rpURL)
		},
	})
	defer ReverseProxy.Close()

	response, err := http.Get(ReverseProxy.URL)
	if err != nil {
		t.Fatalf("failed to get response from ReverseProxy.URL")
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body")
	}

	responseString := string(b)

	if responseString != "test message\n" {
		t.Errorf("response body string mismatch: expected: %q, actual: %q", "test message\n", responseString)
	}

}

func TestGateHandler_BlockOnce(t *testing.T) {

	serverURL := url.URL{
		Scheme: "http",
		Host:   "localhost:0",
	}

	TB := algorithm.NewTokenBucket(0.0, 0.0, 0.0)
	RP := NewReverseProxy(&serverURL)
	handler := GateHandler(TB, RP)

	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	responseHeader := recorder.Header()
	responseStatus := recorder.Code

	if responseHeader.Get("Retry-After") != "" {
		t.Errorf("expected nil header, got %v", responseHeader)
	}

	if responseStatus != http.StatusTooManyRequests {
		t.Errorf("expected status: TooManyRequests, got %v", responseStatus)
	}

}
