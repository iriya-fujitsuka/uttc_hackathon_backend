package dao

import (
	"database/sql"
	"log"
	"uttc_hackathon_backend/models"
)

func GetPostByID(postID string) (*models.Post, error) {
	var post models.Post
	var replyToID sql.NullString
	query := "SELECT id, user_id, community_id, content, reply_to_id, created_at FROM posts WHERE id = ?"
	err := db.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.CommunityID, &post.Content, &replyToID, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 投稿が見つからない場合
		}
		log.Printf("Error fetching post by ID: %v", err)
		return nil, err
	}
	if replyToID.Valid {
		post.ReplyToID = replyToID.String
	} else {
		post.ReplyToID = "" // NULLの場合は空文字列を設定
	}
	return &post, nil
}
