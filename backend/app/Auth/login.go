package auth

import (
	"encoding/json"
	"fmt"
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
	err := db.DB.QueryRow(`SELECT pr.id,pr.display_name, pe.password_hash  FROM profile pr
	 	JOIN user pe ON pr.id = pe.id
	 	WHERE pe.email = ? OR pr.display_name = ?;`, credentials.NoE, credentials.NoE).Scan(&id, &userName, &storedPassword)

	if err != nil || !credentials.Password.Verify([]byte(storedPassword)) {
		logs.ErrorLog.Println("Login failed for user:", credentials.NoE)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "Invalid username or password"}`)
		return
	}
	fmt.Println("Login successful for user:", userName, "ID:", id)
	authorize(w, r, id)
	JsRespond(w, "Login successful", 200)
	
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"message": "Login successful",
	// })
}
