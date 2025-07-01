package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	db "social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

// Login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials structs.Login

	json.NewDecoder(r.Body).Decode(&credentials)

	var storedPassword string
	var id int
	var userName string
	err := db.DB.QueryRow(`
	SELECT pr.id, pr.display_name, pe.password_hash  
	FROM profile pr
	JOIN user pe ON pr.id = pe.id
	WHERE LOWER(pe.email) = LOWER(?) OR LOWER(pr.display_name) = LOWER(?)
`, credentials.NoE, credentials.NoE).Scan(&id, &userName, &storedPassword)

	if err != nil || !auth.CheckPassword(credentials.Password, id) {
		logs.ErrorLog.Println("Error logging in:", err)
		auth.JsRespond(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	auth.Authorize(w, r, id)
	auth.JsRespond(w, "Login successful", 200)
}
