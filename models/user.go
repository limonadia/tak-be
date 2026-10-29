package models

import (
	"time"
)

type User struct {
	ID int `json:"id"`
	Email string  `json:"email"`
	PasswordHash string `json:"password"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	City string `json:"city"`
	Birthday time.Time `json:"birthday"`
	ProfilePicture string `json:"profilePicture"`
	Description string `json:"description"`
	Language Language `json:"language"`
	Theme Theme `json:"theme"`
	CreatedAt time.Time `json:"createdAt"`
	StrikesCount int `json:"strikesCount"`
	Preferences []uint64 `json:"preferences"`
	ResetToken string `json:"-"`
	ResetTokenExpiry time.Time `json:"-"`
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