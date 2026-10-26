package main

import (
	"fmt"
	"log"
	"tak-be/database"
	"tak-be/routes"
	"github.com/gorilla/mux"
	"os"
	"net/http"
	"github.com/rs/cors"
)

func main() {
	db, err := database.Connect()

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer db.Close()

	fmt.Print("Successfully connected to tak!")

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8100", },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(":"+port, handler))

}