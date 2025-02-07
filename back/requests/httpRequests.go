package requests

import (
	"back/db"
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
func UpdateContainers(w http.ResponseWriter, r *http.Request) {
	var containers []db.Container
	if err := json.NewDecoder(r.Body).Decode(&containers); err != nil {
		http.Error(w, "Invalid input!", http.StatusBadRequest)
		return
	}
	result := db.UpdateContainers(&containers)
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
