package main

import (
	"back/db"
	"back/requests"
	"net/http"

	"github.com/rs/cors"
)

const (
	PORT = ":3010"
)

func main() {
	db.InitDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			requests.GetContainers(w, r)
		case http.MethodPost:
			requests.UpdateContainers(w, r)
		default:
			http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
		}
	})
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://pinger", "http://front", "http://localhost:5137"}, // Разрешенные домены
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	http.ListenAndServe(PORT, c.Handler(mux))
}
