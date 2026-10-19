package routes

import (
	"tak-be/handlers"
	"net/http"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/v1").Subrouter()
}