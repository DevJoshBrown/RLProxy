package proxy

import (
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

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

func GateHandler(TB *algorithm.TokenBucket, RP *httputil.ReverseProxy, Metrics *metrics) http.Handler {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			result, retry := TB.Allow()

			if result {
				RP.ServeHTTP(w, r)
				Metrics.requestsTotal.WithLabelValues("allowed").Inc()
			} else {
				if retry >= 0 {
					seconds := math.Ceil(retry.Seconds())
					intsec := int(seconds)
					strsec := strconv.Itoa(intsec)
					w.Header().Set("Retry-After", strsec)
				}
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				Metrics.requestsTotal.WithLabelValues("denied").Inc()
			}
		})
	return handler
}
