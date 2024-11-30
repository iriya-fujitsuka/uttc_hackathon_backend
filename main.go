package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"uttc_hackathon_backend/config"
	"uttc_hackathon_backend/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// データベース初期化
	config.InitDB()
	defer config.CloseDB()

	// 環境変数からポートを取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	// ルーター設定
	router := mux.NewRouter()

	// ユーザー関連ルート
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	router.HandleFunc("/users", handlers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/login", handlers.HandleLogin).Methods("POST")

	// シグナルキャッチでDBを閉じる
	closeDBWithSysCall()

	// サーバー起動
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

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
