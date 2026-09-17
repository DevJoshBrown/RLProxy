package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
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
