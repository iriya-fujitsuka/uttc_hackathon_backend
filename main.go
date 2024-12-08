package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"uttc_hackathon_backend/controller"
	"uttc_hackathon_backend/dao"
)

func CORSMiddlewareProd(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// プリフライトリクエストの応答
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 次のミドルウェアまたはハンドラを呼び出す
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting server...")
	// Initialize database connection
	dao.InitDB()

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users", controller.UserHandler)
	router.HandleFunc("/api/posts", controller.PostHandler)
	router.HandleFunc("/api/replies", controller.ReplyHandler)
	router.HandleFunc("/api/toggle-like", controller.ToggleLike)
	router.HandleFunc("/api/like-counts", controller.GetLikeCount)
	router.HandleFunc("/api/users", controller.GetUserByEmailHandler)
	router.HandleFunc("/api/communities", controller.GetCommunitiesHandler).Methods("GET")
	router.HandleFunc("/api/communities", controller.AddCommunityHandler).Methods("POST")
	router.HandleFunc("/api/communities", controller.DeleteCommunityHandler).Methods("DELETE")
	
	// Wrap the router with the CORS middleware
	handler := CORSMiddlewareProd(router)

	// Handle system call for graceful shutdown
	// handleSysCall()
	log.Println("before env")
	// 環境変数からポートを取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}
	log.Printf("after env")
	// Start HTTP server
	log.Printf("Listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// func handleSysCall() {
// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
// 	go func() {
// 		s := <-sig
// 		log.Printf("Received signal: %v", s)
// 		dao.CloseDB()
// 		log.Println("Database connection closed.")
// 		os.Exit(0)
// 	}()
// }
