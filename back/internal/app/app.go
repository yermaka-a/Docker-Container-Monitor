package app

import (
	"back/internal/config"
	"back/internal/db"
	"back/internal/services"

	"fmt"
	"net/http"

	"github.com/rs/cors"
)

const (
	PORT = ":3010"
)

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			services.GetContainers(w, r)
		case http.MethodPost:
			services.UpdateContainers(w, r)
		default:
			http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
		}
	})
	db.InitDB()
	conn, ch := services.ConnContainersRMQ()
	defer conn.Close()
	defer ch.Close()
	FRONT_PORT := config.New().FrontConfig.FRONT_PORT
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://pinger", "http://front", fmt.Sprintf("http://localhost:%s", FRONT_PORT)}, // Разрешенные домены
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	http.ListenAndServe(PORT, c.Handler(mux))
}
