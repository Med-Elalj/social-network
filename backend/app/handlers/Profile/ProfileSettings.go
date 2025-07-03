package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/modules"
	"social-network/server/logs"
)

func ProfileSettingsHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Extract JWT user payload from context
	payload := r.Context().Value(auth.UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if !ok {
		http.Error(w, "Unauthorized - invalid user context", http.StatusUnauthorized)
		return
	}

	// Step 2: Match request by "type" in path
	switch r.PathValue("type") {
	case "updateProfile":
		UpdateProfileSettings(w, r, data.Sub)
	case "updatePassword":
		UpdatePasswordSettings(w, r, data.Sub)
	case "updateEmail":
		UpdateEmailSettings(w, r, data.Sub)
	case "delete":
		DeleteUserProfile(w, r, data.Sub)
	default:
		http.Error(w, "Invalid action type", http.StatusBadRequest)
	}
}

func UpdateProfileSettings(w http.ResponseWriter, r *http.Request, userID int) {
	var body struct {
		Nickname string `json:"nickname" validate:"required,min=2,max=50"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.Nickname == "" {
		http.Error(w, "Missing nickname", http.StatusBadRequest)
		return
	}

	// Do DB update using userID
	// Example: db.Exec("UPDATE users SET nickname = ? WHERE id = ?", body.Nickname, userID)
	if _, err := modules.DB.Exec("UPDATE profile SET display_name = ? WHERE id = ?", body.Nickname, userID); err != nil {
		logs.ErrorLog.Println("Failed to update nickname:", err)
		auth.JsRespond(w, "Failed to update nickname", http.StatusInternalServerError)
		return
	}
	auth.JsRespond(w, "Nickname updated", http.StatusOK)
}

func UpdateEmailSettings(w http.ResponseWriter, r *http.Request, userID int) {
	// This function should handle email, not nickname
	var body struct {
		Email string `json:"email"` // Not nickname!
	}
	// Validation and response messages should also reference email

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.Email == "" {
		http.Error(w, "Missing email", http.StatusBadRequest)
		return
	}

	if _, err := modules.DB.Exec("UPDATE profile SET email = ? WHERE id = ?", body.Email, userID); err != nil {
		logs.ErrorLog.Println("Failed to update email:", err)
		auth.JsRespond(w, "Failed to update email", http.StatusInternalServerError)
		return
	}

	auth.JsRespond(w, "Nickname updated", http.StatusOK)
}

func UpdatePasswordSettings(w http.ResponseWriter, r *http.Request, userID int) {
	var body struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if body.CurrentPassword == "" || body.NewPassword == "" {
		http.Error(w, "Both current and new passwords are required", http.StatusBadRequest)
		return
	}

	if len(body.NewPassword) < 8 {
		http.Error(w, "New password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Example: check password, hash new password, update in DB
	// user, _ := db.GetUserByID(userID)
	// if !CheckPasswordHash(body.CurrentPassword, user.PasswordHash) { ... }

	auth.JsRespond(w, "Password updated", http.StatusOK)
}

func DeleteUserProfile(w http.ResponseWriter, r *http.Request, userID int) {
	var body struct {
		ConfirmDelete bool `json:"confirmDelete"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Confirm deletion logic
	// db.Exec("DELETE FROM users WHERE id = ?", userID)

	auth.JsRespond(w, "User deleted", http.StatusOK)
}
