package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/DevJoshBrown/RLProxy/internal/algorithm"
)

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	RP := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
		},
	}
	return RP
}

func GateHandler(TB *algorithm.TokenBucket, RP *httputil.ReverseProxy) http.Handler {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if TB.Allow() {
				RP.ServeHTTP(w, r)
			} else {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			}
		})
	return handler
}
