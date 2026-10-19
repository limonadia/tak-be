package models

import (
	"encoding/json"
	"time"
)

type Chat struct {
	ID int `json:"id"`
	Type Type `json:"type"`
	EventId int `json:"event_id"`
	GroupId int `json:"group_id"`
	UserOneId int `json:"user_one_id"`
	UserTwoId int `json:"user_two_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID int `json:"id"`
	ChatId int `json:"chat_id"`
	SenderId int `json:"sender_id"`
	MessageText string `json:"message_text"`
	CreatedAt time.Time `json:"created_at"`
}

const Type (
	DIRECT string = 'direct'
	EVENT string = 'event'
	GROUP string = 'group'
)