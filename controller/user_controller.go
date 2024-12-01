package controller

import (
	"log"
	"net/http"

	"uttc_hackathon_backend/usecase"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		usecase.HandleUserSearch(w, r)
	case http.MethodPost:
		usecase.HandleUserRegister(w, r)
	default:
		log.Printf("Unsupported HTTP Method: %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}
