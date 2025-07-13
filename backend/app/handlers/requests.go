package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "social-network/app/Auth"
	logs "social-network/app/logs"
	"social-network/app/modules"
)

func GetRequestsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var tpdefined int
	if err := json.NewDecoder(r.Body).Decode(&tpdefined); err != nil {
		logs.ErrorLog.Printf("Failed to decode request body: %q", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest) // error
		return
	}

	requests, err := modules.GetRequests(uid, tpdefined)
	if err != nil {
		logs.ErrorLog.Println("Error getting requests:", err)
		auth.JsRespond(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	fmt.Println(requests)

	json.NewEncoder(w).Encode(requests)
}
