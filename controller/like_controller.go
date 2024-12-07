package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"uttc_hackathon_backend/dao"
)

func ToggleLike(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID string `json:"user_id"`
		PostID string `json:"post_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == "" || request.PostID == "" {
		http.Error(w, "User ID and Post ID are required", http.StatusBadRequest)
		return
	}

	if err := dao.ToggleLike(request.UserID, request.PostID); err != nil {
		log.Printf("Error toggling like: %v", err)
		http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetLikeCount(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("postID")
	if postID == "" {
		http.Error(w, `{"error": "Post ID is required"}`, http.StatusBadRequest)
		return
	}

	count, err := dao.GetLikeCount(postID)
	if err != nil {
		log.Printf("Error getting like count: %v", err)
		http.Error(w, `{"error": "Failed to get like count"}`, http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf(`{"like_count": %d}`, count)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}
