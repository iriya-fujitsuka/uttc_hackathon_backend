package usecase

import (
	"encoding/json"
	"log"
	"net/http"

	"uttc_hackathon_backend/dao"
)

func HandleUserSearch(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	users, err := dao.GetUserByName(name)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
