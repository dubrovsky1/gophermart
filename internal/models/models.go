package models

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID
	Login    string
	Password string
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
