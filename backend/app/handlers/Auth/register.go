package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/structs"
	"social-network/server/logs"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user auth.Register

	json.NewDecoder(r.Body).Decode(&user)
	if len(user.UserName) == 0 {
		user.UserName = auth.GenerateNickname(user.Fname, user.Lname)
		if user.UserName == "" {
			auth.JsRespond(w, "Please enter a valid username.", http.StatusBadRequest)
		}
	}
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

	userID, err := auth.InsertUser(user)
	if err != nil {
		logs.ErrorLog.Println("Error inserting user into database:", err)
		if structs.SqlConstraint(&err) {
			auth.JsRespond(w, "Username or email already exists", http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error": "sorry couldn't register your information please try aggain later"}`)
			return
		}
	}
	auth.Authorize(w, r, int(userID))
	auth.JsRespond(w, "Registration successful", http.StatusOK)
}
