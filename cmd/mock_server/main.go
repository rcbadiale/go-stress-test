package main

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Register struct {
	Requests int
	mu       *sync.Mutex
}

func main() {
	register := Register{Requests: 0, mu: &sync.Mutex{}}
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		prob := rand.Float32()
		status := http.StatusOK
		if prob < 0.2 {
			status = http.StatusInternalServerError
		} else if prob < 0.4 {
			status = http.StatusServiceUnavailable
		} else if prob < 0.6 {
			status = http.StatusNotFound
		} else if prob < 0.8 {
			status = http.StatusUnauthorized
		}
		// Throttle the server for testing purposes
		time.Sleep(time.Duration(1+prob) * time.Second)
		w.WriteHeader(status)

		register.mu.Lock()
		defer register.mu.Unlock()
		register.Requests += 1
		log.Println("Request", register.Requests, "Status", status)
	})
	http.ListenAndServe(":8080", nil)
}
