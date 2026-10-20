package models

import (
	
)

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ImageUrl string `json:"image_url"`
}