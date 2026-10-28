package handlers

import (
	"tak-be/database"
	"tak-be/models"
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	"strings"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User

	err := database.DB.QueryRow(`
		SELECT id, email, first_name, last_name, city, birthday,
		       profile_picture, description
		FROM users
		WHERE id = ?
	`, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.City,
		&user.Birthday,
		&user.ProfilePicture,
		&user.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Println("Database error:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("Failed to encode profile:", err)
	}
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updates map[string]json.RawMessage

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	allowedFields := map[string]string{
		"email":           "email",
		"first_name":      "first_name",
		"last_name":       "last_name",
		"city":            "city",
		"birthday":        "birthday",
		"profile_picture": "profile_picture",
		"description":     "description",
		"language":        "language",
		"theme":           "theme",
	}

	var setClauses []string
	var args []interface{}

	for field, column := range allowedFields {
		rawValue, exists := updates[field]
		if !exists {
			continue
		}

		var value interface{}
		if err := json.Unmarshal(rawValue, &value); err != nil {
			http.Error(w, "Invalid value for "+field, http.StatusBadRequest)
			return
		}

		setClauses = append(setClauses, column+"=?")
		args = append(args, value)
	}

	if len(setClauses) == 0 {
		http.Error(w, "No profile fields to update", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET " + strings.Join(setClauses, ", ") + " WHERE id=?"
	args = append(args, userID)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		log.Println("Update error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	GetProfile(w, r)
}

func GetTheme(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var theme models.Theme
	err := database.DB.QueryRow(`
		SELECT theme FROM users WHERE id=?`, userID).Scan(&theme)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to retrieve user theme: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(theme)
}

func GetLanguage(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var language models.Language
	err := database.DB.QueryRow(`
		SELECT language FROM users WHERE id=?`, userID).Scan(&language)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to retrieve user language %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(language)
}

func GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`
		SELECT i.id, i.name
		FROM interests i
		INNER JOIN user_interests ui
			ON ui.interest_id = i.id
		WHERE ui.user_id = ?
		ORDER BY i.name ASC
	`, userID)

	if err != nil {
		log.Println("Failed to fetch preferences:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	preferences := make([]models.Interest, 0)

	for rows.Next() {
		var interest models.Interest

		if err := rows.Scan(&interest.ID, &interest.Name); err != nil {
			log.Println("Failed to scan preference:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		preferences = append(preferences, interest)
	}

	if err := rows.Err(); err != nil {
		log.Println("Failed to read preferences:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(preferences); err != nil {
		log.Println("Failed to encode preferences:", err)
	}
}

func UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		Preferences *[]uint64 `json:"preferences"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Ensure the preferences field was provided.
	// An empty array is allowed and clears all preferences.
	if request.Preferences == nil {
		http.Error(w, "Preferences field is required", http.StatusBadRequest)
		return
	}

	preferences := *request.Preferences

	// Start transaction.
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Remove the user's existing preferences.
	_, err = tx.Exec(
		`DELETE FROM user_interests WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		log.Println("Failed to delete preferences:", err)
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	// Insert the new preferences.
	for _, interestID := range preferences {
		_, err = tx.Exec(
			`INSERT INTO user_interests (user_id, interest_id)
			 VALUES (?, ?)`,
			userID,
			interestID,
		)
		if err != nil {
			log.Println("Failed to insert preference:", err)
			http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
			return
		}
	}

	// Commit changes only when everything succeeds.
	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit preferences:", err)
		http.Error(w, "Failed to save preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Preferences updated successfully",
		"preferences": preferences,
	})
}