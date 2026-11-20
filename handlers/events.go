package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tak-be/database"
	"tak-be/models"
	"time"
)

func CreateEvent(w http.ResponseWriter, r *http.Request){
	userID, ok := GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 6<<20)

	if err := r.ParseMultipartForm(6<<20); err != nil {
		log.Println("MULTIPART ERROR:", err)
		http.Error(w, "Invalid multipart request", http.StatusBadRequest)
		return
	}

	var req models.CreateEventRequest
	if err := json.Unmarshal([]byte(r.FormValue("event")), &req); err != nil {
		log.Println("JSON ERROR:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	startsAt, err := time.Parse(
        "2006-01-02 15:04:05",
        req.Date+" "+req.TimeFrom,
    )
    if err != nil {
        log.Println("START TIME ERROR:", err)
        http.Error(w, "Invalid start date or time", http.StatusBadRequest)
        return
    }

	endsAt, err := time.Parse(
        "2006-01-02 15:04:05",
        req.Date+" "+req.TimeTo,
    )
    if err != nil {
        log.Println("END TIME ERROR:", err)
        http.Error(w, "Invalid end date or time", http.StatusBadRequest)
        return
    }

	registrationDeadline, err := time.Parse(
        "2006-01-02T15:04:05",
        req.RegistrationDeadline,
    )
    if err != nil {
        log.Println("DEADLINE ERROR:", err)
        http.Error(w, "Invalid registration deadline", http.StatusBadRequest)
        return
    }

	imagePath, err := saveImage(r, "image", "events")
	if err != nil {
		log.Println("EVENT IMAGE ERROR:", err)
		http.Error(w, "Invalid event image", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec(
        `INSERT INTO events (
            creator_id,
            title,
            description,
            image_url,
            location_name,
            city,
            address,
            latitude,
            longitude,
            starts_at,
            ends_at,
            registration_deadline,
            min_participants,
            max_participants,
            price_amount,
            currency,
            status
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        userID,
        req.Title,
        req.Description,
        imagePath,
        req.LocationName,
        req.City,
        req.Address,
        req.Latitude,
        req.Longitude,
        startsAt,
        endsAt,
        registrationDeadline,
        req.MinParticipants,
        req.MaxParticipants,
        req.PriceAmount,
        req.Currency,
        models.PUBLISHED,
    )
	if err != nil {
		log.Println("DATABASE ERROR:", err)
		http.Error(w, "Could not create event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("LAST INSERT ID ERROR:", err)
		http.Error(w, "Could not retrieve event id", http.StatusInternalServerError)
		return
	}

	for _, interestID := range req.Labels {
		_, err := database.DB.Exec(
			`INSERT INTO event_interests (event_id, interest_id) VALUES (?, ?)`, id, interestID,
		)
		if err != nil {
			log.Printf("Failed to insert interest %d for user %d: %v", interestID, id, err)
			http.Error(w, "Could not save user interests", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	 json.NewEncoder(w).Encode(map[string]interface{}{
        "id":          id,
        "creatorId":   userID,
        "title":       req.Title,
        "description": req.Description,
        "image":       imagePath,
        "location":    req.LocationName,
        "city":        req.City,
        "address":     req.Address,
        "latitude":    req.Latitude,
        "longitude":   req.Longitude,
        "startsAt":    startsAt,
        "endsAt":      endsAt,
        "labels":      req.Labels,
        "relatedGroups": req.RelatedGroups,
        "status":      models.PUBLISHED,
    })
}