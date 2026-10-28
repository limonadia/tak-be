package handlers

import (
	"tak-be/database"
	"tak-be/models"
	"net/http"
	"encoding/json"
	"log"
)

func GetInterests(w http.ResponseWriter, r *http.Request) {
	if database.DB == nil {
		http.Error(w, "database unavailable", http.StatusInternalServerError)
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, name FROM interests ORDER BY name ASC`,
	)
	if err != nil {
		http.Error(w, "Failder to fetch interests", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    interests := make([]models.Interest, 0)

    for rows.Next() {
        var interest models.Interest

        if err := rows.Scan(&interest.ID, &interest.Name); err != nil {
			http.Error(w, "Failed to read interests", http.StatusInternalServerError)
            return
        }

        interests = append(interests, interest)
    }

    if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read interests", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(interests); err != nil {
		log.Printf("Failed to encode interests %v", err)
	}
}