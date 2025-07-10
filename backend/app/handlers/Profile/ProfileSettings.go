package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/logs"
	"social-network/app/modules"
)

func ProfileSettingsHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Extract JWT user payload from context
	payload := r.Context().Value(auth.UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if !ok {
		auth.JsRespond(w, "Unauthorized - invalid user context", http.StatusUnauthorized)
		return
	}

	fmt.Println(r.PathValue("type"))
	// Step 2: Match request by "type" in path
	switch r.PathValue("type") {
	case "updateUsername":
		UpdateUsername(w, r, data.Sub)
	case "updatePassword":
		UpdatePassword(w, r, data.Sub)
	case "delete":
		DeleteProfile(w, r, data.Sub)
	case "changePrivacy":
		ChangePrivacy(w, r, data.Sub)
	default:
		auth.JsRespond(w, "Invalid action type", http.StatusBadRequest)
	}
}

func ChangePrivacy(w http.ResponseWriter, r *http.Request, userID int) {
	// Define request structure
	var body struct {
		Privacy bool `json:"privacy"`
	}

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logs.ErrorLog.Println("JSON decode error:", err)
		auth.JsRespond(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate input
	if r.Body == http.NoBody {
		auth.JsRespond(w, "Request body is required", http.StatusBadRequest)
		return
	}

	// Update database
	_, err := modules.DB.Exec(
		"UPDATE profile SET is_public = ? WHERE id = ?",
		body.Privacy,
		userID,
	)
	if err != nil {
		logs.ErrorLog.Println("Database update error:", err)
		auth.JsRespond(w, "Failed to update privacy setting", http.StatusInternalServerError)
		return
	}
	// Success response
	auth.JsRespond(w, "Privacy setting updated", http.StatusOK)
	logs.InfoLog.Printf("Privacy setting updated to %v for user ID: %d", body.Privacy, userID)
}

func UpdateUsername(w http.ResponseWriter, r *http.Request, userID int) {
	var body struct {
		Nickname string `json:"nickname" validate:"required,min=2,max=50"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		auth.JsRespond(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.Nickname == "" {
		auth.JsRespond(w, "Missing nickname", http.StatusBadRequest)
		return
	}
	if !auth.IsValidNickname(body.Nickname) {
		auth.JsRespond(w, "Nickname must be 3-13 characters and use letters or underscores.", http.StatusBadRequest)
		return
	}
	if auth.NicknameExists(body.Nickname) {
		auth.JsRespond(w, "Nickname already exists", http.StatusConflict)
		return
	}

	if _, err := modules.DB.Exec("UPDATE profile SET display_name = ? WHERE id = ?", body.Nickname, userID); err != nil {
		logs.ErrorLog.Println("Failed to update nickname:", err)
		auth.JsRespond(w, "Failed to update nickname", http.StatusInternalServerError)
		return
	}
	auth.JsRespond(w, "Nickname updated", http.StatusOK)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request, userID int) {
	var body struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	if body.CurrentPassword == "" || body.NewPassword == "" {
		auth.JsRespond(w, "Both current and new passwords are required", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		auth.JsRespond(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if !(auth.IsValidPassword(body.NewPassword) && auth.IsValidPassword(body.ConfirmPassword) && body.NewPassword != body.ConfirmPassword) {
		auth.JsRespond(w, "New Password is not valid", http.StatusBadRequest)
		return
	}
	if !auth.CheckPassword(body.CurrentPassword, userID) {
		auth.JsRespond(w, "Current password is incorrect", http.StatusBadRequest)
		return
	}
	auth.ChangePassword(body.NewPassword, int64(userID))
	auth.JsRespond(w, "Password updated", http.StatusOK)
}

func DeleteProfile(w http.ResponseWriter, r *http.Request, userID int) {
	var body struct {
		ConfirmDelete  bool   `json:"confirmDelete"`
		DeletePassword string `json:"deletePassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		auth.JsRespond(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if body.DeletePassword == "" {
		auth.JsRespond(w, "Delete password is required", http.StatusBadRequest)
		return
	}
	if !auth.CheckPassword(body.DeletePassword, userID) {
		auth.JsRespond(w, "Delete password is incorrect", http.StatusBadRequest)
		return
	}
	// Confirm deletion
	if !body.ConfirmDelete {
		auth.JsRespond(w, "Please confirm deletion", http.StatusBadRequest)
		return
	}
	// Delete user profile from the database
	if _, err := modules.DB.Exec("DELETE FROM profile WHERE id = ?", userID); err != nil {
		logs.ErrorLog.Println("Failed to delete user profile:", err)
		auth.JsRespond(w, "Failed to delete user profile", http.StatusInternalServerError)
		return
	}
	logs.InfoLog.Println("User profile deleted successfully for user ID:", userID)
	auth.JsRespond(w, "User deleted", http.StatusOK)
}
