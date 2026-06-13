package models

import "time"

// User is an account. PasswordHash is tagged json:"-" so it NEVER leaks into
// an API response (Session 15's struct tags protecting a secret).
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// RegisterInput is the body for POST /auth/register.
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginInput is the body for POST /auth/login.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
