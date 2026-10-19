package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID int `json:"id"`
	Email string  `json:"email"`
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

const Language (
	EN string = 'en'
	AL string = 'al'
)

const Theme (
	DARK string = 'dark'
	LIGHT string = 'light'
)