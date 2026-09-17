package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/DevJoshBrown/RLProxy/internal/algorithm"
	"github.com/DevJoshBrown/RLProxy/internal/proxy"
)

func main() {
	serverURL := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}

	TB := algorithm.NewTokenBucket(1000, 25, 1000)
	Proxy := proxy.NewReverseProxy(&serverURL)

	gateHandler := proxy.GateHandler(TB, Proxy)
	log.Fatal(http.ListenAndServe(":8081", gateHandler))
}
