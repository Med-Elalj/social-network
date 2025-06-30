package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

func GroupCreation(w http.ResponseWriter, r *http.Request, uid int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logs.ErrorLog.Println("Error reading request body:", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Request body cannot be empty"}`)
		return
	}

	var group structs.Group

	json.NewDecoder(r.Body).Decode(&group)

	// validating what in group creation
	// if err := group.Validate(); err != nil {
	// 	logs.Println("Validation failed for group:", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, `{"error": %q}`, err.Error())
	// 	return
	// }
	// group.Cid = structs.ID(uid)
	_, err = modules.InsertGroup(group, uid)
	if err != nil {
		auth.JsRespond(w, "group creation failed", http.StatusInternalServerError)
	}
	auth.JsRespond(w, "group posted successfully", http.StatusOK)
}
