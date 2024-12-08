package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc_hackathon_backend/dao"
	"uttc_hackathon_backend/models"
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

func HandlePostReply(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 必要なフィールドが空でないことを確認
	if post.UserID == "" || post.ReplyToID == "" {
		http.Error(w, "User ID and ReplyToID are required", http.StatusBadRequest)
		return
	}

	// 元の投稿のコミュニティIDを取得
	originalPost, err := dao.GetPostByID(post.ReplyToID)
	if err != nil {
		log.Printf("Error fetching original post: %v", err)
		http.Error(w, "Failed to fetch original post", http.StatusInternalServerError)
		return
	}
	post.CommunityID = originalPost.CommunityID

	postID, err := dao.AddPost(post)
	if err != nil {
		log.Printf("Error adding reply: %v", err)
		http.Error(w, "Failed to create reply", http.StatusInternalServerError)
		return
	}
	post.ID = postID

	w.WriteHeader(http.StatusCreated)
}
