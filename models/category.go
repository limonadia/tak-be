package models

import (
	"encoding/json"
)

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ImageUrl string `json:"image_url"`
}