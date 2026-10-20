package models

import (
	"time"
)

type Connection struct {
	ID int `json:"id"`
	RequesterId int `json:"requester_id"`
	RecipientId int `json:"recipient_id"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Status string
const (
	PENDING Status = "pending"
	ACCEPTED Status = "accepted"
	REJECTED Status = "rejected"
)