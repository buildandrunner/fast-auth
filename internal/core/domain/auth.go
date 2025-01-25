package domain

import "time"

type User struct {
	ID          string `json:"id"`
	Phonenumber string `json:"phonenumber"`
}

type Session struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
