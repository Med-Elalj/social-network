package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	var user structs.Register

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
// 		return
// 	}

// 	if validationErrors := user.ValidateRegister(); len(validationErrors) != 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"error": validationErrors,
// 		})
// 		return
// 	}

// 	if err := db.InsertUser(user); err != nil {
// 		logs.ErrorLog.Println("DB error inserting user:", err)
// 		if structs.SqlConstraint(&err) {
// 			w.WriteHeader(http.StatusConflict)
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"error": err.Error(),
// 			})
// 		} else {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"error": "Registration failed. Try again later.",
// 			})
// 		}
// 		return
// 	}

// 	// Get inserted user ID from DB safely (optional: RETURNING id)
// 	userID := GetElemVal("id", "user", `display_name = "`+string(user.UserName)+`"`).(int)
// 	fmt.Println("Registered user ID:", userID)

// 	// authorize(w, r, userID)

//		// Write response after cookies are set
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(map[string]interface{}{
//			"message": "Registration successful",
//		})
//	}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user structs.Register

	json.NewDecoder(r.Body).Decode(&user)

	if err := user.ValidateRegister(); len(err) != 0 {
		logs.ErrorLog.Println("Validation failed for user input")
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
		logs.ErrorLog.Println("Error inserting user into database:", err)
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
