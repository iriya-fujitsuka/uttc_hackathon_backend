package usecase

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"uttc_hackathon_backend/dao"

	"github.com/oklog/ulid/v2"
)

type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	var userReq UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if userReq.Name == "" || userReq.Email == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	if err := dao.AddUser(id, userReq.Name, userReq.Email); err != nil {
		log.Printf("Error inserting user: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
