package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var requestCount = 0

type res struct {
	Counter int       `json:"counter"`
	Time    time.Time `json:"time"`
}

func incrementalCounter(count int) (newCount int) {
	count++
	return count
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	requestCount = incrementalCounter(requestCount)

	resp := res{
		Counter: requestCount,
		Time:    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

func main() {
	http.HandleFunc("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
