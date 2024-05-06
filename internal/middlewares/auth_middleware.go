package middlewares

import (
	"log"
	"net/http"
	"os"
)

func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Auth")

		if apiKey == "" || apiKey != os.Getenv("api_key") {
			log.Printf("Unauthorized API Token: %v", apiKey)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte("{error: \"Unauthorized\"}"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
