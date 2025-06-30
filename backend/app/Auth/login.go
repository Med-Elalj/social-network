package auth

import (
	"encoding/json"
	"net/http"

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

	if err != nil || !CheckPassword(credentials.Password, id) {
		logs.ErrorLog.Println("Error logging in:", err)
		JsRespond(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	authorize(w, r, id)
	JsRespond(w, "Login successful", 200)
}
