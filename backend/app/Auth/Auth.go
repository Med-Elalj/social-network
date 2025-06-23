package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social-network/app/Auth/jwt"
	db "social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

type contextKey string

const UserContextKey contextKey = "user"

// Check if user logged
func Islogged(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := r.Context().Value(UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "Unauthorized"}`)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": data.Username + " User is logged in",
	})
}

// Login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials structs.Login

	json.NewDecoder(r.Body).Decode(&credentials)

	var storedPassword string
	var id int
	var userName string
	err := db.DB.QueryRow(`SELECT pr.id,pr.display_name, pe.password_hash  FROM profile pr
	 	JOIN person pe ON pr.id = pe.id
	 	WHERE pe.email = ? OR pr.display_name = ?;`, credentials.NoE, credentials.NoE).Scan(&id, &userName, &storedPassword)

	if err != nil || !credentials.Password.Verify([]byte(storedPassword)) {
		logs.Println("Login failed for user:", credentials.NoE)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "Invalid username or password"}`)
		return
	}

	cookie := &http.Cookie{
		Name:     "JWT",
		Value:    jwt.Generate(jwt.CreateJwtPayload(id, userName)),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(jwt.Time_to_Expire),
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
	})
}

// Logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:     "JWT",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,                   // ensures deletion
		HttpOnly: true,                 // match original
		Secure:   true,                 // match original
		SameSite: http.SameSiteLaxMode, // match original
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

// Register handler
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user structs.Register

	json.NewDecoder(r.Body).Decode(&user)

	if err := user.ValidateRegister(); len(err) != 0 {
		logs.Println("Validation failed for user input")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}
	// if user.Avatar != "" {
	// upload.UploadHandler(w, r, user.Avatar)
	// 	// TODO:Uploading avatar
	// }

	err := db.InsertUser(user)
	if err != nil {
		fmt.Println(err)
		logs.Println("Error inserting user into database:", err)
		if structs.SqlConstraint(&err) {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, `{"error": %q}`, err.Error())
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error": "sorry couldn't register your information please try aggain later"}`)
			return
		}
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
	})
}
