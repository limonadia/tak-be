package handlers

import (
	"tak-be/database"
	"tak-be/models"
)

var jwtKey = []byte("password")

func Signup(w http.ResponseWriter, r *http.Request){
	var user models.User
}