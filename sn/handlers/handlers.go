package handlers

import (
	"fmt"
	"io"
	"net/http"

	"social-network/server/logs"
	"social-network/sn/db"
	"social-network/sn/structs"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	// TODO: ADD middle ware to check if the user is already logged in
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

	// TODO hash password
	user.Password.Hash()
	n, err := db.DB.Exec(`
    BEGIN TRANSACTION;
INSERT INTO profile (display_name, is_person) VALUES (?, 1);
INSERT INTO person (ent, email, first_name, last_name, password_hash, date_of_birth,gender)
VALUES (last_insert_rowid(), ?, ?, ?, ?, ?, ?);

COMMIT;
`,
		user.UserName, user.Email, user.Fname, user.Lname, user.Password, user.Birthdate, user.Gender)
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
	// var credentials struct {
	// 	Username string `json:"username"`
	// 	Password string `json:"password"`
	// }
	// if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
	// 	http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
	// 	return
	// }

	// var storedPassword string
	// var id string
	// var usrnm string
	// err := db.DB.QueryRow("SELECT password, id FROM users WHERE username = ?", credentials.Username).Scan(&storedPassword, &id)
	// if err != nil || storedPassword != credentials.Password {
	// 	err = db.DB.QueryRow("SELECT password, id FROM users WHERE email = ?", credentials.Username).Scan(&storedPassword, &id)
	// 	if err != nil || storedPassword != credentials.Password {
	// 		http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
	// 		return
	// 	} else {
	// 		err = db.DB.QueryRow("SELECT username  FROM users WHERE email = ?", credentials.Username).Scan(&usrnm)
	// 		if err != nil {
	// 			http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
	// 			return
	// 		}
	// 		credentials.Username = usrnm
	// 	}
	// }

	// err = db.DB.QueryRow("SELECT password, id FROM users WHERE username = ? OR email = ?", credentials.Username, credentials.Username).Scan(&storedPassword, &id)
	// if err != nil || storedPassword != credentials.Password {
	// 	http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
	// 	return
	// }

	// jwt.Generate()

	// // Set cookie
	// cookie := &http.Cookie{
	// 	Name:     "JWT",
	// 	Value:    "jwt",
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	Expires:  time.Now().Add(168 * time.Hour),
	// }

	// http.SetCookie(w, cookie)
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message":  "Login successful",
	// 	"username": credentials.Username,
	// })
}

func Forunf(w http.ResponseWriter, r *http.Request) {
	// TODO implement this function
}
