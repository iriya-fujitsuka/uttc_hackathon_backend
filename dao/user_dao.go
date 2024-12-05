package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"uttc_hackathon_backend/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitDB() {
	// .envファイルの読み込み
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
	// DB接続のための準備
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	// connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPwd, mysqlHost, "3306", mysqlDatabase)
	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	log.Printf("Connecting to database: %s\n", connStr)
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	// if err := db.Ping(); err != nil {
	// 	log.Fatalf("Database is unreachable: %v\n", err)
	// }
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v\n", err)
	}
}

func GetUserByName(name string) ([]models.User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var id, name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			return nil, err
		}
		users = append(users, models.User{Id: id, Name: name, Email: email})
	}
	return users, nil
}

func AddUser(id, name, email string) error {
	query := "INSERT INTO users (id, name, email) VALUES (?, ?, ?)"
	_, err := db.Exec(query, id, name, email)
	if err != nil {
		log.Printf("Database insertion failed: %v\n", err)
	}
	return err
}

func AddPost(post models.Post) error {
	query := "INSERT INTO posts (user_id, content, reply_to_id) VALUES (?, ?, ?)"
	_, err := db.Exec(query, post.UserID, post.Content, post.ReplyToID)
	if err != nil {
		log.Printf("Failed to insert post: %v\n", err)
		return err
	}
	return nil
}

// 投稿を取得 (すべての投稿を取得)
func GetPosts() ([]models.Post, error) {
	rows, err := db.Query("SELECT id, user_id, content, reply_to_id, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ReplyToID, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// 特定投稿の返信を取得
func GetReplies(postID string) ([]models.Post, error) {
	rows, err := db.Query("SELECT id, user_id, content, reply_to_id, created_at FROM posts WHERE reply_to_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ReplyToID, &post.CreatedAt); err != nil {
			return nil, err
		}
		replies = append(replies, post)
	}
	return replies, nil
}