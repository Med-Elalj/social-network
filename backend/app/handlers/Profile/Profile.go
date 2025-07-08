package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

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
	IsPublic    bool           `json:"isPublic"`
	IsUser      bool           `json:"isUser"`
	CreatedAt   string         `json:"created_at"`
	IsSelf      bool           `json:"isSelf"`
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	nickname := strings.TrimSpace(r.PathValue("name"))
	payload, ok := r.Context().Value(auth.UserContextKey).(*jwt.JwtPayload)

	var profile Profile
	var err error

	// 👤 Case 1: viewer requests their own profile using their nickname
	if ok && strings.EqualFold(nickname, payload.Username) {
		// Fetch by ID (self profile)
		err = modules.DB.QueryRow(`
			SELECT email, first_name, last_name, display_name, date_of_birth, gender,
			       avatar, description, is_public, is_user, created_at
			FROM profile WHERE id = ?
		`, payload.Sub).Scan(
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
		profile.IsSelf = true
	} else {
		// 🕵️‍♂️ Case 2: someone else’s profile
		err = modules.DB.QueryRow(`
			SELECT email, first_name, last_name, display_name, date_of_birth, gender,
			       avatar, description, is_public, is_user, created_at
			FROM profile
			WHERE LOWER(display_name) = LOWER(?)
		`, nickname).Scan(
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

		if profile.IsUser && !profile.IsPublic {
			profile.Email = ""
			profile.FirstName = ""
			profile.LastName = ""
			profile.DateOfBirth = ""
			profile.Gender = ""
		}

		profile.IsSelf = false
	}

	if err != nil {
		if err == sql.ErrNoRows {
			auth.JsRespond(w, "Profile not found", http.StatusNotFound)
		} else {
			logs.ErrorLog.Println("DB error:", err)
			auth.JsRespond(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
