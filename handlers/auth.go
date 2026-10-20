package handlers

import (
	"tak-be/database"
	"tak-be/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"io"
	"database/sql"
	"fmt"
)

var jwtKey = []byte("password_hash")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func Signup(w http.ResponseWriter, r *http.Request){

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("Signup JSON decoding error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.PasswordHash == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var exists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email=?", user.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	if database.DB == nil {
		log.Println("Database is not initialized")
		http.Error(w, "Database unavailable", http.StatusInternalServerError)
		return
	}

	user.Language = models.EN
	user.Theme = models.DARK
	user.StrikesCount = 0
	res, err := database.DB.Exec(
		"INSERT INTO users (email, password_hash, first_name, last_name, city, birthday, profile_picture, description, language, theme, strikes_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Email, string(hashed), user.FirstName, user.LastName, user.City, user.Birthday, user.ProfilePicture, user.Description, user.Language, user.Theme, user.StrikesCount, 
	)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	user.ID = int(id)
	user.PasswordHash = ""
	log.Printf("Signup inserted user ID=%d, email=%s", id, user.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email string `json:"email"`
		PasswordHash string `json:"password_hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	row := database.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", creds.Email)
	if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
		log.Println("Login failed: user not found for email: ", creds.Email)
		http.Error(w, "Email does not exist", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.PasswordHash)); err != nil {
		log.Println("Password comparison failed: ", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || token == nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, email, first_name, last_name, city, birthday, profile_picture, description, theme, language
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.City, &user.Birthday, &user.ProfilePicture, &user.Description
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
	json.NewEncoder(w).Encode(user)
}

func GetTheme(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var theme models.Theme
	err := database.DB.QueryRow(`
		SELECT theme FROM users WHERE id = ?`, userID).Scan(&theme)

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
		SELECT language FROM users WHERE id = ?`, userID).Scan(&language)
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