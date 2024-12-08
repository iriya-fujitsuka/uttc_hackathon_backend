package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"uttc_hackathon_backend/dao"
)

func GetCommunitiesHandler(w http.ResponseWriter, r *http.Request) {
	communities, err := dao.GetCommunities()
	if err != nil {
		log.Printf("Error fetching communities: %v", err)
		http.Error(w, "Failed to fetch communities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(communities); err != nil {
		log.Printf("Error encoding communities: %v", err)
		http.Error(w, "Failed to encode communities", http.StatusInternalServerError)
	}
}

func AddCommunityHandler(w http.ResponseWriter, r *http.Request) {
	var community struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&community); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := dao.AddCommunity(community.Name)
	if err != nil {
		log.Printf("Error adding community: %v", err)
		http.Error(w, "Failed to add community", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func DeleteCommunityHandler(w http.ResponseWriter, r *http.Request) {
	var community struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&community); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := dao.DeleteCommunity(community.ID); err != nil {
		log.Printf("Error deleting community: %v", err)
		http.Error(w, "Failed to delete community", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"id": %d, "status": "deleted"}`, community.ID)))
}