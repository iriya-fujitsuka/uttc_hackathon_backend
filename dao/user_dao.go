package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	// DB接続のための準備
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v\n", err)
	}
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v\n", err)
	}
}

func GetUserByName(name string) ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, name, age FROM user WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id, name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			return nil, err
		}
		users = append(users, map[string]interface{}{
			"id":   id,
			"name": name,
			"age":  age,
		})
	}
	return users, nil
}

func AddUser(id, name string, age int) error {
	_, err := db.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", id, name, age)
	return err
}
