package main

import (
	"github.com/gorilla/mux"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"uttc_hackathon_backend/config"
	"uttc_hackathon_backend/handlers"
)

func main() {
	// データベース初期化
	config.InitDB()
	defer config.CloseDB()

	// ルーター設定
	router := mux.NewRouter()

	// ユーザー関連ルート
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	router.HandleFunc("/users", handlers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/login", handlers.HandleLogin).Methods("POST")

	// 環境変数からポートを取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	// サーバー停止シグナルをキャッチしてDBをクローズ
	closeDBWithSysCall()

	// サーバー起動
	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

// シグナルキャッチでDBを閉じる
func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("Received signal: %v", s)
		config.CloseDB()
		os.Exit(0)
	}()
}
