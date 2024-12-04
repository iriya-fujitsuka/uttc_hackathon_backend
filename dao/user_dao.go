package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	model "uttc_hackathon_backend/models"
	userController "uttc_hackathon_backend/controller/user"
	postController "uttc_hackathon_backend/controller/post"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitDB() {
	// .envファイルの読み込み
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// DB接続のための準備
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPwd, mysqlHost, "3306", mysqlDatabase)
	// 	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
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

func GetUserByName(name string) ([]model.User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var id, name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			return nil, err
		}
		users = append(users, model.User{Id: id, Name: name, Email: email})
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

func AddPost(content string) error {
	query := "INSERT INTO posts (content) VALUES (?)"
	_, err := db.Exec(query, content)
	if err != nil {
		log.Printf("Failed to insert post: %v\n", err)
		return err
	}
	return nil
}