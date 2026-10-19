package models

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID int `json:"id"`
	Recipient_id int `json:"recipient_id"`
	ActorId int `json:"actor_id"`
	Type Type `json:"type"`
	EntityType EntityType `json:"entity_type"`
	EntityId int `json:"entity_id"`
	Title string `json:"title"`
	Body string `json:"body"`
	IsRead boolean `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt time.Time `json:"read_at"`
}

const Type (
	CONNECTION_REQUEST string = 'connection_request'
	CONNECTION_ACCEPTED string = 'connection_accepted'
	EVENT_INVITATION string = 'event_invitation'
	EVENT_UPDATE string = 'event_update'
	EVENT_CANCELLED string = 'event_cancelled'
	EVENT_REMINDER string = 'event_reminder'
	GROUP_INVITATION string = 'group_invitation'
	GROUP_UPDATE string = 'group_update'
	NEW_MESSAGE string = 'new_message'
	EVENT_CHAT_MESSAGE string = 'event_chat_message'
	GROUP_CHAT_MESSAGE string = 'group_chat_message'
	USER_REPORT_UPDATE string = 'user_report_update'
	STRIKE_RECIEVED string = 'strike_received'
)

const EntityType (
	USER string = 'user'
	EVENT string = 'event'
	GROUP string = 'group'
	CHAT string = 'chat'
	REPORT string = 'report'
	STRIKE string = 'strike'
)