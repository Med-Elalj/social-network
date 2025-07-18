package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
)

// handler
func GetSearchHandler(w http.ResponseWriter, r *http.Request, uid int) {
	query := r.URL.Query().Get("query")
	offsetst := r.URL.Query().Get("offset")
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		groupId = 0
	}
	offset, err := strconv.Atoi(offsetst)
	if err != nil {
		logs.ErrorLog.Printf("Invalid offset: %q", err)
		auth.JsResponse(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	profiles, err := modules.GetSearchprofile(query, offset, groupId, uid)
	if err != nil {
		logs.ErrorLog.Printf("Error getting search profiles: %q", err)
		auth.JsResponse(w, "Failed to get search profiles", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profiles)
}
