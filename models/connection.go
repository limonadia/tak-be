package models

import (
	"encoding/json"
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

const Status (
	PENDING string = 'pending'
	ACCEPTED string = 'accepted'
	REJECTED string = 'rejected'
)