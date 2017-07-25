package web

import (
	"log"
	"net/http"
	"time"
)



func TimeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startingTime := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Now().Sub(startingTime)
		log.Printf("request took %v\n", duration)
	})
}
