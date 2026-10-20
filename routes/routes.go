package routes

import (
	"tak-be/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/v1").Subrouter()
	
	api.HandleFunc("/signup", handlers.Signup).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")

	api.Handle("/getTheme", handlers.AuthMiddleWare(http.HandlerFunc(handlers.GetTheme))).Methods("GET")
	api.Handle("/getLanguage", handlers.AuthMiddleWare(http.HandlerFunc(handlers.GetLanguage))).Methods("GET")
	api.Handle("/getProfile", handlers.AuthMiddleWare(http.HandlerFunc(handlers.GetProfile))).Methods("GET")
}