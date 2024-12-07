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

func AddPost(post models.Post) (string, error) {
	var query string
	var err error
	var result sql.Result

	if post.ReplyToID == "" {
		query = "INSERT INTO posts (user_id, content, reply_to_id) VALUES (?, ?, NULL)"
		result, err = db.Exec(query, post.UserID, post.Content)
	} else {
		query = "INSERT INTO posts (user_id, content, reply_to_id) VALUES (?, ?, ?)"
		result, err = db.Exec(query, post.UserID, post.Content, post.ReplyToID)
	}

	if err != nil {
		log.Printf("Failed to insert post: %v\n", err)
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v\n", err)
		return "", err
	}

	return fmt.Sprintf("%d", id), nil // IDを文字列として返す
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
		var replyToID sql.NullString
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &replyToID, &post.CreatedAt); err != nil {
			return nil, err
		}
		if replyToID.Valid {
			post.ReplyToID = replyToID.String
		} else {
			post.ReplyToID = ""
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

func GetRepliesByPostID(postID string) ([]models.Post, error) {
	rows, err := db.Query("SELECT id, user_id, content, reply_to_id, created_at FROM posts WHERE reply_to_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []models.Post
	for rows.Next() {
		var post models.Post
		var replyToID sql.NullString
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &replyToID, &post.CreatedAt); err != nil {
			return nil, err
		}
		if replyToID.Valid {
			post.ReplyToID = replyToID.String
		} else {
			post.ReplyToID = ""
		}
		replies = append(replies, post)
	}
	return replies, nil
}

func ToggleLike(userID, postID string) error {
	log.Printf("Toggling like for user %s and post %s", userID, postID)

	// ユーザーが存在するか確認
	var userExists bool
	userCheckQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	err := db.QueryRow(userCheckQuery, userID).Scan(&userExists)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return err
	}
	if !userExists {
		log.Printf("User %s does not exist", userID)
		return fmt.Errorf("user does not exist")
	}

	// すでに「いいね」されているかを確認
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)"
	err = db.QueryRow(checkQuery, userID, postID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking like existence: %v", err)
		return err
	}

	if exists {
		// すでに「いいね」されている場合は削除
		deleteQuery := "DELETE FROM likes WHERE user_id = ? AND post_id = ?"
		_, err = db.Exec(deleteQuery, userID, postID)
		if err != nil {
			log.Printf("Error deleting like: %v", err)
			return err
		}
		log.Printf("Like removed for user %s and post %s", userID, postID)
	} else {
		// まだ「いいね」されていない場合は追加
		insertQuery := "INSERT INTO likes (user_id, post_id) VALUES (?, ?)"
		_, err = db.Exec(insertQuery, userID, postID)
		if err != nil {
			log.Printf("Error adding like: %v", err)
			return err
		}
		log.Printf("Like added for user %s and post %s", userID, postID)
	}

	return nil
}

// GetLikeCount は特定の投稿に対する「いいね」の数を取得します。
func GetLikeCount(postID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM likes WHERE post_id = ?"
	err := db.QueryRow(query, postID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // いいねがない場合は0を返す
		}
		log.Printf("Error getting like count: %v", err)
		return 0, err
	}
	return count, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ユーザーが見つからない場合
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userID string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email FROM users WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ユーザーが見つからない場合
		}
		return nil, err
	}
	return &user, nil
}

func AddReply(postID, replyContent string) error {
	query := "INSERT INTO replies (post_id, content) VALUES (?, ?)"
	_, err := db.Exec(query, postID, replyContent)
	if err != nil {
		log.Printf("Error adding reply: %v", err)
		return err
	}
	return nil
}
