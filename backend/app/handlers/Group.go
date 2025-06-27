package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

// anas
func GroupJoin(w http.ResponseWriter, r *http.Request, uid int) {
	var group structs.GroupReq

	json.NewDecoder(r.Body).Decode(&group)

	// validating what in group creation
	// if err := group.Validate(); err != nil {
	//     logs.Println("Validation failed for group:", err)
	//     w.WriteHeader(http.StatusBadRequest)
	//     fmt.Fprintf(w, `{"error": %q}`, err.Error())
	//     return
	// }
	// group.Uid = uid
	if uid != group.Uid {
		///request
	}
	if !modules.InsertUGroup(group, uid) {
		structs.JsRespond(w, "group joining failed", http.StatusInternalServerError)
	}

	// TODO notif to group creator
	structs.JsRespond(w, "user send req group successfully", http.StatusOK)
}

func GroupCreation(w http.ResponseWriter, r *http.Request, uid int) {
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

	var group structs.Group

	json.NewDecoder(r.Body).Decode(&group)

	// validating what in group creation
	// if err := group.Validate(); err != nil {
	// 	logs.Println("Validation failed for group:", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, `{"error": %q}`, err.Error())
	// 	return
	// }
	group.Cid = structs.ID(uid)
	_, err = modules.InsertGroup(group)
	if err != nil {
		structs.JsRespond(w, "group creation failed", http.StatusInternalServerError)
	}
	structs.JsRespond(w, "group posted successfully", http.StatusOK)
}
