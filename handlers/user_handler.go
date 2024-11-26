package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"uttc_hackathon_backend/config"
	"uttc_hackathon_backend/models"
	"uttc_hackathon_backend/utils"
)

// 新規ユーザー作成
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	user.ID = utils.GenerateULID()

	query := "INSERT INTO users (id, firebase_uid, name, email) VALUES (?, ?, ?, ?)"
	_, err := config.DB.Exec(query, user.ID, user.Name, user.Email)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// 全ユーザー取得
func GetUsers(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, firebase_uid, name, email, deleted_at FROM users WHERE deleted_at IS NULL"
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.DeletedAt); err != nil {
			http.Error(w, "Failed to parse users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// ユーザー論理削除
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	firebaseUID := r.Header.Get("FirebaseUID")
	query := "UPDATE users SET deleted_at = ? WHERE firebase_uid = ?"
	_, err := config.DB.Exec(query, time.Now(), firebaseUID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User logically deleted successfully"}`))
}

// ユーザーログイン
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// リクエストボディのデコード
	var loginReq struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// データベースでメールアドレスを検索
	var user models.User
	query := "SELECT id, email FROM users WHERE email = ?"
	err := config.DB.QueryRow(query, loginReq.Email).Scan(&user.ID, &user.Email)
	if err != nil {
		log.Printf("Login failed: User not found for email %s", loginReq.Email)
		http.Error(w, "Invalid email address", http.StatusUnauthorized)
		return
	}

	// ログイン成功レスポンス
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Login successful",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
