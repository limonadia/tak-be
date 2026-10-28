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
	"fmt"
	"github.com/google/uuid"
	"path/filepath"
	"os"
	"encoding/hex"
	"math/rand"
	"net/smtp"
)

var jwtKey = []byte("password_hash")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func Signup(w http.ResponseWriter, r *http.Request){
	r.Body = http.MaxBytesReader(w, r.Body, 6<<20)

	if err := r.ParseMultipartForm(6<<20); err != nil {
		http.Error(w, "Invalid multipart request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.Unmarshal([]byte(r.FormValue("user")), &user); err != nil {
		log.Printf("Signup JSON decoding error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.PasswordHash == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email=?", user.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error:" +err.Error(), http.StatusInternalServerError)
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

	profilePath, err := saveProfilePicture(r)
	if err != nil {
		http.Error(w, "Invalid profile picture", http.StatusBadRequest)
		return
	}

	user.ProfilePicture = profilePath
	user.Language = models.EN
	user.Theme = models.DARK
	user.StrikesCount = 0

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	res, err := database.DB.Exec(
		"INSERT INTO users (email, password_hash, first_name, last_name, city, birthday, profile_picture, description, language, theme, strikes_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Email, string(hashed), user.FirstName, user.LastName, user.City, user.Birthday, user.ProfilePicture, user.Description, user.Language, user.Theme, user.StrikesCount,
	)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Could not retrieve userID", http.StatusInternalServerError)
		return
	}

	user.ID = int(id)

	for _, interestID := range user.Preferences {
		_, err := tx.Exec(
			`INSERT INTO user_interests (user_id, interest_id) VALUES (?, ?)`, id, interestID,
		)
		if err != nil {
			log.Printf("Failed to insert interest %d for user %d: %v", interestID, id, err)
			http.Error(w, "Could not save user interests", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Could not complete signup", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""
	log.Printf("Signup inserted user ID=%d, email=%s", id, user.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func saveProfilePicture(r *http.Request) (string, error) {
    file, header, err := r.FormFile("profile_picture")
    if err != nil {
        if err == http.ErrMissingFile {
            return "", nil 
        }
        return "", err
    }
    defer file.Close()

    if header.Size > 5<<20 {
        return "", fmt.Errorf("image exceeds 5 MB")
    }

    buffer := make([]byte, 512)
    n, err := file.Read(buffer)
    if err != nil && err != io.EOF {
        return "", err
    }

    contentType := http.DetectContentType(buffer[:n])

    extensions := map[string]string{
        "image/jpeg": ".jpg",
        "image/png":  ".png",
        "image/webp": ".webp",
    }

    ext, ok := extensions[contentType]
    if !ok {
        return "", fmt.Errorf("unsupported image format")
    }

    filename := uuid.NewString() + ext
    relativePath := filepath.ToSlash(filepath.Join("uploads", "profiles", filename))
    fullPath := filepath.Join("uploads", "profiles", filename)

    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return "", err
    }

    dst, err := os.Create(fullPath)
    if err != nil {
        return "", err
    }

    if _, err := file.Seek(0, io.SeekStart); err != nil {
        dst.Close()
        os.Remove(fullPath)
        return "", err
    }

    _, copyErr := io.Copy(dst, file)
    closeErr := dst.Close()

    if copyErr != nil {
        os.Remove(fullPath)
        return "", copyErr
    }
    if closeErr != nil {
        os.Remove(fullPath)
        return "", closeErr
    }

    return relativePath, nil
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

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Token == "" || body.NewPassword == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var userID int 
	var expiry time.Time
	err := database.DB.QueryRow("SELECT id, reset_token_expiry FROM users WHERE reset_token=?", body.Token).Scan(&userID, &expiry)
	if err != nil {
		http.Error(w, "Token expired", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("UPDATE users SET password=?, reset_token=NULL, reset_token_expiry=NULL WHERE id =?", string(hashed), userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successful.",
	})
}

func generateRandomToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("Error generating random token: ", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userID)
	if err != nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	resetToken := hex.EncodeToString(tokenBytes)

	_, err = database.DB.Exec(
		"UPDATE users SET reset_token = ?, reset_token_expiry = ? WHERE id = ?",
		resetToken, time.Now().Add(1*time.Hour), userID,
	)
	if err != nil {
		http.Error(w, "Failed to store reset token", http.StatusInternalServerError)
		return
	}

	resetLink := "https://localhost:8100/reset-password?token=" + resetToken
	fmt.Println("Password reset link: ", resetLink)
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASS")
	to := req.Email
	subject := "Password Reset Request"
	body := fmt.Sprintf("You requested to reset your password. \n\nClick the link below: \n%s\n\nThis link will expire in 1 hour.", resetLink)
	message := []byte("Subject: " + subject + "\r\n" +
	"From: " + from + "\r\n" +
	"To: " + to + "\r\n\r\n" +
	body)
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message); err != nil {
		fmt.Println("Email error:", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset link sent successfully"})
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction: ", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM user_interests WHERE user_id = ?`, userID)
	if err != nil {
		log.Println("Failed to delete user interests: ", err)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec(`DELETE FROM users WHERE id = ?`, userID)
	if err != nil {
		log.Println("Failed to delete user:", err)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to check deleted account", http.StatusInternalServerError)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit deletion:", err)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{ "message": "Account deleted successfully" })
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}