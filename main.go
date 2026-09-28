package main

import (
	"fmt"
	"log"
	"tak-be/database"
)

func main() {
	db, err := database.Connect()

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer db.Close()

	fmt.Print("Successfully connected to tak db!")

	//Add APIs later
}