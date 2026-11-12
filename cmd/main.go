package main

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/DevJoshBrown/RLProxy/internal/algorithm"
	"github.com/DevJoshBrown/RLProxy/internal/proxy"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	upstreamScheme := os.Getenv("UPSTREAMSCHEME")
	if upstreamScheme == "" {
		slog.Error("no defined upstream server scheme", "error", upstreamScheme)
		os.Exit(1)
	}

	upstreamHost := os.Getenv("UPSTREAMHOST")
	if upstreamHost == "" {
		slog.Error("no defined upsteam server host", "error", upstreamHost)
		os.Exit(1)
	}

	upstreamURL := url.URL{
		Scheme: upstreamScheme,
		Host:   upstreamHost,
	}

	TB := algorithm.NewTokenBucket(1000, 25, 1000)
	Proxy := proxy.NewReverseProxy(&upstreamURL)

	reg := prometheus.NewRegistry()
	Metrics := proxy.NewMetrics(reg)

	gateHandler := proxy.GateHandler(TB, Proxy, Metrics)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.Handle("/", gateHandler)

	proxyPort := os.Getenv("PROXYPORT")
	if proxyPort == "" {
		proxyPort = "8081"
	}

	slog.Info("proxy server starting", "port", proxyPort)

	proxyServer := &http.Server{
		Addr:         ":" + proxyPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	err := proxyServer.ListenAndServe()
	if err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("ok"))
}
