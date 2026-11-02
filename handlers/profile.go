package handlers

import (
	"tak-be/database"
	"tak-be/models"
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	"strings"
	"os"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	var profilePicture sql.NullString

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
		&profilePicture,
		&user.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("GET PROFILE DATABASE ERROR: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	} else {
		user.ProfilePicture = ""
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

	r.Body = http.MaxBytesReader(w, r.Body, 6<<20)

	if err := r.ParseMultipartForm(6 << 20); err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	userJSON := r.FormValue("user")
	if userJSON == "" {
		http.Error(w, "Missing user data", http.StatusBadRequest)
		return
	}

	var updates map[string]json.RawMessage

	if err := json.Unmarshal([]byte(userJSON), &updates); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	var removeProfilePicture bool

	if raw, exists := updates["removeProfilePicture"]; exists {
		if err := json.Unmarshal(raw, &removeProfilePicture); err != nil {
			http.Error(w, "Invalid removeProfilePicture value", http.StatusBadRequest)
			return
		}
	}

	var preferences *[]uint64
	
	if rawPreferences, exists := updates["preferences"]; exists {
		var parsed []uint64
		if err := json.Unmarshal(rawPreferences, &parsed); err != nil {
			http.Error(w, "Invalid preferences", http.StatusBadRequest)
			return
		}
		preferences = &parsed
	}

	newProfilePath, err := saveProfilePicture(r)
	if err != nil {
		http.Error(w, "Invalid profile picture: "+err.Error(), http.StatusBadRequest)
		return
	}
	if newProfilePath != "" {
		updates["profilePicture"], _ = json.Marshal(newProfilePath)
	} else if removeProfilePicture {
		updates["profilePicture"], _ = json.Marshal(nil)
	}

	allowedFields := map[string]string{
		"email":           "email",
		"firstname":      "first_name",
		"lastname":       "last_name",
		"city":            "city",
		"birthday":        "birthday",
		"profilePicture": "profile_picture",
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

	hasPreferences := preferences != nil

	if len(setClauses) == 0 && !hasPreferences{
		if newProfilePath != "" {
            os.Remove(newProfilePath)
        }
		http.Error(w, "No profile fields to update", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
if err != nil {
    log.Printf("Failed to begin profile update transaction: %v", err)

    if newProfilePath != "" {
        os.Remove(newProfilePath)
    }

    http.Error(w, "Database error", http.StatusInternalServerError)
    return
}
	defer tx.Rollback()

	if len(setClauses) > 0 {
		query := "UPDATE users SET " + strings.Join(setClauses, ", ") + " WHERE id=?"
		args = append(args, userID)
		_, err = tx.Exec(query, args...)
		if err != nil {
			if newProfilePath != "" {
                os.Remove(newProfilePath)
            }
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}
	}

	if preferences != nil {
		if err := updatePreferencesTx(tx, userID, *preferences); err != nil {
			if newProfilePath != "" {
                os.Remove(newProfilePath)
            }
			http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		if newProfilePath != "" {
			os.Remove(newProfilePath)
		}
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
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
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	preferences := make([]models.Interest, 0)

	for rows.Next() {
		var interest models.Interest

		if err := rows.Scan(&interest.ID, &interest.Name); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		preferences = append(preferences, interest)
	}

	if err := rows.Err(); err != nil {
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

	// Start transaction.
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if err := updatePreferencesTx(tx, userID, *request.Preferences); err != nil {
		http.Error(w, "Failed to update preferences:", http.StatusInternalServerError)
		return
	}

	// Commit changes only when everything succeeds.
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to save preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Preferences updated successfully",
		"preferences": *request.Preferences,
	})
}

func updatePreferencesTx(tx *sql.Tx, userID int, preferences []uint64) error {
	_, err := tx.Exec(
		`DELETE FROM user_interests WHERE user_id = ?`, userID,
	)
	if err != nil {
		return err
	}

	for _, interestID := range preferences {
		_, err := tx.Exec(
			`INSERT INTO user_interests (user_id, interest_id) VALUES (?, ?)`, userID, interestID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}