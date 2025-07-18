package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	db "social-network/app/modules"
	"social-network/app/structs"
)

// Login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials structs.Login

	json.NewDecoder(r.Body).Decode(&credentials)

	var storedPassword string
	var id int
	var userName string
	var avatar sql.NullString
	err := db.DB.QueryRow(`
  SELECT pr.id, pr.display_name, pe.password_hash, pr.avatar
  FROM profile pr
  JOIN user pe ON pr.id = pe.id
  WHERE LOWER(pr.email) = LOWER(?) OR LOWER(pr.display_name) = LOWER(?)
`, credentials.NoE, credentials.NoE).
		Scan(&id, &userName, &storedPassword, &avatar)

	if err != nil || !auth.CheckPassword(credentials.Password, id) {
		logs.ErrorLog.Println("Error logging in:", err)
		auth.JsResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	auth.Authorize(w, r, id)
	auth.JsMapResponse(w, map[string]any{
		"message":  "Login successful",
		"avatar":   avatar.String,
	}, http.StatusOK)
}
