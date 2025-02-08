package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pinger/internal/models"
)

func PostData(containersL []models.Container) {
	const url = "http://web:80/containers"
	jsonData, err := json.Marshal(containersL)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)
}
