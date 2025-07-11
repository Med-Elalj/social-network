package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
)

func GetRequestsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var tpdefined int
	json.NewDecoder(r.Body).Decode(&tpdefined)
	requests, err := modules.GetRequests(uid, tpdefined)
	if err != nil {
		logs.ErrorLog.Println("Error getting requests:", err)
		auth.JsRespond(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}
