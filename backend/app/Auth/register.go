package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

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
