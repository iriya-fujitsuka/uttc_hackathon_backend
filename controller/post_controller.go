package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc_hackathon_backend/models"
	"uttc_hackathon_backend/usecase"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Received request: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

	switch r.Method {
	case http.MethodPost:
		log.Println("Handling POST request")
		usecase.HandlePostCreate(w, r)
	case http.MethodGet:
		log.Println("Handling GET request")
		usecase.HandlePostFetch(w, r)
	default:
		log.Printf("Unsupported HTTP Method: %s\n", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}