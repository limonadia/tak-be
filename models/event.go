package models 

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID int `json:"id"`
	CreatorId int `json:"creator_id"`
	CategoryId int `json:"category_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	ImageUrl string `json:"image_url"`
	LocationName string `json:"location_name"`
	City string `json:"city"`
	Address string `json:"address"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt time.Time `json:"ends_at"`
	RegistrationDeadline time.Time `json:"registration_deadline"`
	MinParticipants int `json:"min_participants"`
	MaxParticipants int `json:"max_participants"`
	PriceAmount float64 `json:"price_amount"`
	Currency string `json:"currency"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const Status (
	DRAFT string = 'draft'
	PUBLISHED string = 'published'
	CANCELLED string = 'cancelled'
)