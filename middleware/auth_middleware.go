package middleware

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"

	"firebase.google.com/go/v4/auth"
)

var firebaseAuth *auth.Client

func init() {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	firebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}
}

// トークン検証ミドルウェア
func VerifyFirebaseToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		idToken := authHeader[7:]
		token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Set("FirebaseUID", token.UID)
		next.ServeHTTP(w, r)
	}
}
