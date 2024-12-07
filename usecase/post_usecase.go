package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc_hackathon_backend/dao"
	"uttc_hackathon_backend/models"
)

func HandlePostCreate(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// user_idが空でないことを確認
	if post.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if err := dao.AddPost(post); err != nil {
		log.Printf("Error adding post: %v", err) // エラーログを追加
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandlePostList(w http.ResponseWriter, r *http.Request) {
	posts, err := dao.GetPosts()
	if err != nil {
		log.Printf("Error fetching posts: %v", err) // エラーログを追加
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Error encoding posts: %v", err) // エラーログを追加
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
	}
}
