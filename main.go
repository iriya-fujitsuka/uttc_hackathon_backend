package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"uttc_hackathon_backend/controller"
	"uttc_hackathon_backend/dao"
)

func main() {
	// Initialize database connection
	dao.InitDB()

	// Set up routes
	http.HandleFunc("/user", controller.UserHandler)

	// Handle system call for graceful shutdown
	handleSysCall()

	// 環境変数からポートを取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	// Start HTTP server
	log.Printf("Listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("Received signal: %v", s)
		dao.CloseDB()
		log.Println("Database connection closed.")
		os.Exit(0)
	}()
}
