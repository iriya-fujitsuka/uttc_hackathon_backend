package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"uttc_hackathon_backend/dao"
	"uttc_hackathon_backend/usecase"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // JSONレスポンスを明示
	log.Printf("Received request: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

	switch r.Method {
	case http.MethodGet:
		log.Println("Handling GET request")
		usecase.HandleUserSearch(w, r)
	case http.MethodPost:
		log.Println("Handling POST request")
		usecase.HandleUserRegister(w, r)
	default:
		log.Printf("Unsupported HTTP Method: %s\n", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// ユーザーIDをURLから取得
	userID := r.URL.Path[len("/api/users/"):]

	// データベースからユーザー情報を取得
	user, err := dao.GetUserByID(userID)
	if err != nil {
		http.Error(w, "ユーザー情報の取得に失敗しました。", http.StatusInternalServerError)
		return
	}

	// ユーザー情報をJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "ユーザー情報のエンコードに失敗しました。", http.StatusInternalServerError)
	}
}

func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータからemailを取得
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// データベースからユーザー情報を取得
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "ユーザー情報の取得に失敗しました。", http.StatusInternalServerError)
		return
	}

	// ユーザー情報をJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "ユーザー情報のエンコードに失敗しました。", http.StatusInternalServerError)
	}
}
