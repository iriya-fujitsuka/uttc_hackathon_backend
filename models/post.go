package models

type Post struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	CommunityID string `json:"community_id"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
}