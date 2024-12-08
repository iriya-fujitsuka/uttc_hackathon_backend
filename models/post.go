package models

type Post struct {
    ID            string `json:"id"`             // Postの一意のID (ULIDやUUIDを使用可能)
	UserID        string `json:"user_id"`        // 投稿者のID
	CommunityID   int    `json:"community_id"`   // 投稿が属するコミュニティID
	CommunityName string `json:"community_name"` // 投稿が属するコミュニティ名
	Content       string `json:"content"`        // 投稿の内容
	CreatedAt     string `json:"created_at"`     // 作成日時
	ReplyToID     string `json:"reply_to_id,omitempty"` // 返信先のPost ID。NULLの場合は通常の投稿
}
