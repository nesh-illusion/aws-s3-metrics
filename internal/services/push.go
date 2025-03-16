package services

import (
	"bytes"
	"log"
	"net/http"
)

func PushData(endpoint string, data []byte) {
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Printf("Data pushed successfully. Status: %s", resp.Status)
}
