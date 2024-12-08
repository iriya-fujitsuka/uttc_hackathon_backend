package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"uttc_hackathon_backend/dao"
	"uttc_hackathon_backend/models"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location  = "asia-northeast1"
	modelName = "gemini-1.5-flash-002"
	projectID = "term6-iriya-fujitsuka"
)

func logJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}
	log.Println(string(data))
}

func HandlePostCreate(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Printf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received post: %+v", post) // 受信した投稿データをログに記録

	// user_idとcommunity_idが空でないことを確認
	if post.UserID == "" || post.CommunityID == 0 {
		log.Printf("Missing user_id or community_id: user_id=%s, community_id=%d", post.UserID, post.CommunityID)
		http.Error(w, "User ID and Community ID are required", http.StatusBadRequest)
		return
	}

	postID, err := dao.AddPost(post)
	if err != nil {
		log.Printf("Error adding post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	post.ID = postID

	// 投稿内容をGeminiでチェック
	go checkPostContent(post)

	w.WriteHeader(http.StatusCreated)
}

func checkPostContent(post models.Post) error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Printf("Error creating Gemini client: %v", err)
		return err
	}

	gemini := client.GenerativeModel(modelName)
	promptText := fmt.Sprintf("以下の内容が医学的に問題があるかどうかを判断してください。問題があれば、その理由を述べて警告文を出してください。問題がなければ、内容に共感し寄りそう返信を書いてください。１４０字以内に収めてください: %s", post.Content)
	prompt := genai.Text(promptText)

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	logJSON(map[string]interface{}{
		"level": "info",
		"resp":  *resp.Candidates[0],
	})

	// Candidates配列の中のContent.Partsを使用
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			text, ok := part.(genai.Text)
			if !ok {
				log.Printf("part is not genai.Text")
				continue
			}
			// 新しい投稿を追加し、投稿IDを受け取る
			newPostID, err := dao.AddPost(models.Post{UserID: post.UserID, Content: string(text), ReplyToID: post.ID})
			if err != nil {
				log.Printf("Error adding post: %v", err)
				continue
			}
			log.Printf("Added post: %v, %v", newPostID, string(text))
		}
	}
	return nil
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
