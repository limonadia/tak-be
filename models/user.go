package models

import (
	"time"
)

type User struct {
	ID int `json:"id"`
	Email string  `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	City string `json:"city"`
	Birthday time.Time `json:"birthday"`
	ProfilePicture string `json:"profile_picture"`
	Description string `json:"description"`
	Language Language `json:"language"`
	Theme Theme `json:"theme"`
	CreatedAt time.Time `json:"created_at"`
	StrikesCount int `json:"strikes_count"`
}

type Language string
const (
	EN Language = "en"
	AL Language = "al"
)

type Theme string
const (
	DARK Theme = "dark"
	LIGHT Theme = "light"
)