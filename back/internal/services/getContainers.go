package services

import (
	"back/internal/db"
	"encoding/json"
	"net/http"
)

func GetContainers(w http.ResponseWriter, r *http.Request) {
	containers := db.GetContainers()
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(containers); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
