package models

import(
	"encoding/json"
	"time"
)

type Group struct {
	ID int `json:"id"`
	CreatorId int `json:"creator_id"`
	CategoryId int `json:"category_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	ImageUrl string `json:"image_url"`
	City string `json:"city"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const Status (
	ACTIVE string = 'active'
	ARCHIVED string = 'archived'
)