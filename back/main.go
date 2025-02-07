package main

import (
	"back/db"
	"back/requests"
	"net/http"
)

const (
	PORT = ":3010"
)

func main() {
	db.InitDB()
	http.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			requests.GetContainers(w, r)
		case http.MethodPost:
			requests.UpdateContainers(w, r)
		default:
			http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(PORT, nil)
}
