package controller

import (
	"log"
	"net/http"

	"uttc_hackathon_backend/usecase"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // JSONレスポンスを明示
	log.Printf("Received request: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

	switch r.Method {
	case http.MethodGet:
		log.Println("Handling GET request")
		usecase.HandleUserSearch(w, r)
	case http.MethodPost:
		log.Println("Handling POST request")
		usecase.HandleUserRegister(w, r)
	default:
		log.Printf("Unsupported HTTP Method: %s\n", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
