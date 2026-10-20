package models

import (
	"time"
)

type Notification struct {
	ID int `json:"id"`
	Recipient_id int `json:"recipient_id"`
	ActorId int `json:"actor_id"`
	Type NotificationType `json:"type"`
	EntityType EntityType `json:"entity_type"`
	EntityId int `json:"entity_id"`
	Title string `json:"title"`
	Body string `json:"body"`
	IsRead bool `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt time.Time `json:"read_at"`
}

type NotificationType string
const (
	CONNECTION_REQUEST NotificationType = "connection_request"
	CONNECTION_ACCEPTED NotificationType = "connection_accepted"
	EVENT_INVITATION NotificationType = "event_invitation"
	EVENT_UPDATE NotificationType = "event_update"
	EVENT_CANCELLED NotificationType = "event_cancelled"
	EVENT_REMINDER NotificationType = "event_reminder"
	GROUP_INVITATION NotificationType = "group_invitation"
	GROUP_UPDATE NotificationType = "group_update"
	NEW_MESSAGE NotificationType = "new_message"
	EVENT_CHAT_MESSAGE NotificationType = "event_chat_message"
	GROUP_CHAT_MESSAGE NotificationType = "group_chat_message"
	USER_REPORT_UPDATE NotificationType = "user_report_update"
	STRIKE_RECIEVED NotificationType = "strike_received"
)

type EntityType string
const (
	USER_ENTITY EntityType = "user"
	EVENT_ENTITY EntityType = "event"
	GROUP_ENTITY EntityType = "group"
	CHAT_ENTITY EntityType = "chat"
	REPORT_ENTITY EntityType = "report"
	STRIKE_ENTITY EntityType = "strike"
)