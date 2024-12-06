package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc_hackathon_backend/dao"
)

func HandleGetReplies(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	replies, err := dao.GetRepliesByPostID(postID)
	if err != nil {
		log.Printf("Error fetching replies: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(replies)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
} 