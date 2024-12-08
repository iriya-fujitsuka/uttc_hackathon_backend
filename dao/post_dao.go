package dao

import (
	"database/sql"
	"log"
	"uttc_hackathon_backend/models"
)

func GetPostByID(postID string) (*models.Post, error) {
	var post models.Post
	query := "SELECT id, user_id, community_id, content, reply_to_id, created_at FROM posts WHERE id = ?"
	err := db.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.CommunityID, &post.Content, &post.ReplyToID, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 投稿が見つからない場合
		}
		log.Printf("Error fetching post by ID: %v", err)
		return nil, err
	}
	return &post, nil
} 