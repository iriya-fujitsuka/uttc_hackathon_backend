package models

import "time"

type User struct {
	ID          string     `json:"id"`
	FirebaseUID string     `json:"firebase_uid"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
