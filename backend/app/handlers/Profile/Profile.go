package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/modules"
	"social-network/server/logs"
)

type Profile struct {
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	DisplayName string         `json:"display_name"`
	DateOfBirth string         `json:"date_of_birth,omitempty"`
	Gender      string         `json:"gender"`
	Avatar      sql.NullString `json:"avatar"`
	Description string         `json:"description"`
	IsPublic    bool           `json:"is_public"`
	IsUser      bool           `json:"is_user"`
	CreatedAt   string         `json:"created_at"`
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract Payload from context
	payload, ok := r.Context().Value(auth.UserContextKey).(*jwt.JwtPayload)
	if ok {
		fmt.Printf("User ID: %d, Username: %s\n", payload.Sub, payload.Username)
	}
	uid := payload.Sub

	// 2. Prepare a Profile struct to hold data
	profile := &Profile{}

	// 3. Query the database
	err := modules.DB.QueryRow(`
		SELECT email, first_name, last_name, display_name, date_of_birth, gender,
		       avatar, description, is_public, is_user, created_at
		FROM profile WHERE id = ?
	`, uid).Scan(
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&profile.DisplayName,
		&profile.DateOfBirth,
		&profile.Gender,
		&profile.Avatar,
		&profile.Description,
		&profile.IsPublic,
		&profile.IsUser,
		&profile.CreatedAt,
	)
	fmt.Println("Profile data:", profile)

	// 4. Handle SQL errors
	if err != nil {
		if err == sql.ErrNoRows {
			auth.JsRespond(w, "Profile not found", http.StatusNotFound)
		} else {
			logs.ErrorLog.Println("DB error:", err)
			auth.JsRespond(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// 5. Encode result as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		logs.ErrorLog.Println("Error encoding JSON:", err)
		auth.JsRespond(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
