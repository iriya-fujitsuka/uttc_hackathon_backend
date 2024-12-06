package controller

import (
	"log"
	"net/http"
	"uttc_hackathon_backend/usecase"
)

func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Received request: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

	switch r.Method {
	case http.MethodGet:
		log.Println("Handling GET request for replies")
		usecase.HandleGetReplies(w, r)
	case http.MethodPost:
		log.Println("Handling POST request for reply")
		usecase.HandlePostReply(w, r)
	default:
		log.Printf("Unsupported HTTP Method: %s\n", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
} 