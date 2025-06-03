package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"social-network/server/logs"
	"social-network/sn/db"
	"social-network/sn/security"
	"social-network/sn/security/jwt"
	"social-network/sn/structs"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logs.Println("Error reading request body:", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Request body cannot be empty"}`)
		return
	}

	var user structs.Register

	if err := structs.JsonRestrictedDecoder(body, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": %q}`, err.Error())
		return
	}
	// Validate input

	if err := user.Validate(); err != nil {
		logs.Println("Validation failed for user input")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": %q}`, err.Error())
		return
	}

	n, err := db.InsertUser(user)
	if err != nil {
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

	if user.Avatar != "" {
		// upload.UploadHandler(w, r, user.Avatar)
		// TODO:Uploading avatar
	}

	logs.Println("User inserted into database successfully:", user.UserName)
	id, err := n.LastInsertId()
	logs.Println("User registered successfully with ID:", id)
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error": "sorry couldn't register your information please try aggain later"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"User registered successfully"}`)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logs.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "Request body cannot be empty"}`)
		return
	}
	var credentials structs.Login

	if err := structs.JsonRestrictedDecoder(body, &credentials); err != nil || credentials.Password.IsValid() != nil {
		logs.Errorf("Error decoding request body: %v", err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "Invalid request body"}`)
		return
	}

	var storedPassword string
	var id int
	var userName string
	err = db.DB.QueryRow(
		`SELECT pr.id,pr.display_name, pe.password_hash  FROM profile pr  
JOIN person pe ON pr.id = pe.id
WHERE ? IN (pe.email, pr.display_name);`,
		credentials.NoE).Scan(&id, &userName, &storedPassword)

	if err != nil || !credentials.Password.Verify([]byte(storedPassword)) {
		logs.Println("Login failed for user:", credentials.NoE)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "Invalid username or password"}`)
		return
	}

	// Set cookie
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
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login successful",
		"username": userName,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:    "JWT",
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

func Islogged(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := r.Context().Value(security.UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	fmt.Println(data)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "Unauthorized"}`)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "User is logged in",
		"username": data.Username,
		"id":       data.Sub,
	})
}


