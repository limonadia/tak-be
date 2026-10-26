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
	api.HandleFunc("/interests", handlers.GetIntersts).Methods("GET")

	api.Handle("/theme", handlers.AuthMiddleWare(http.HandlerFunc(handlers.GetTheme))).Methods("GET")
	api.Handle("/language", handlers.AuthMiddleWare(http.HandlerFunc(handlers.GetLanguage))).Methods("GET")
	api.Handle("/profile", handlers.AuthMiddleWare(http.HandlerFunc(handlers.GetProfile))).Methods("GET")
	
}