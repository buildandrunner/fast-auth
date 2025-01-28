package domain

import (
	"time"
)

type User struct {
	ID          string `json:"id"`
	Phonenumber string `json:"phonenumber"`
	Password    string `json:"-"`
}

type Session struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Credentials struct {
	Phonenumber string `json:"phonenumber" form:"phonenumber" binding:"required"`
	Password    string `json:"password" form:"password" binding:"required"`
}
