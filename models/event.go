package models 

import (
	"time"
)

type Event struct {
	ID int `json:"id"`
	CreatorId int `json:"creatorId"`
	Labels []uint64 `json:"labels"`
	Title string `json:"title"`
	Description string `json:"description"`
	ImageUrl string `json:"image"`
	LocationName string `json:"location"`
	City string `json:"city"`
	PlaceId string `json:"placeId"`
	Address string `json:"address"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
	StartsAt time.Time `json:"timeFrom"`
	EndsAt time.Time `json:"timeTo"`
	Date time.Time `json:"date"`
	RegistrationDeadline time.Time `json:"registrationDeadline"`
	MinParticipants int `json:"minPeople"`
	MaxParticipants int `json:"maxPeople"`
	PriceAmount float64 `json:"fee"`
	Currency string `json:"currency"`
	Status EventStatus `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RelatedGroups []int `json:"relatedGroups"`
}

type CreateEventRequest struct {
    Title                string   `json:"title"`
    Description          string   `json:"description"`
    Date                 string   `json:"date"`
    TimeFrom             string   `json:"timeFrom"`
    TimeTo               string   `json:"timeTo"`
    LocationName         string   `json:"location"`
    City                 string   `json:"city"`
    PlaceId              string   `json:"placeId"`
    Address              string   `json:"address"`
    Latitude             float64  `json:"latitude"`
    Longitude            float64  `json:"longitude"`
	CategoryID           uint64   `json:"categoryId"`
    MinParticipants      int      `json:"minPeople"`
    MaxParticipants      int      `json:"maxPeople"`
    RegistrationDeadline string   `json:"registrationDeadline"`
    PriceAmount          float64  `json:"fee"`
    Currency             string   `json:"currency"`
    Labels               []uint64 `json:"labels"`
    RelatedGroups        []int    `json:"relatedGroups"`
}

type EventStatus string
const (
	DRAFT EventStatus = "draft"
	PUBLISHED EventStatus = "published"
	CANCELLED EventStatus = "cancelled"
)