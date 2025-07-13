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
	offset, err := strconv.Atoi(offsetst)
	if err != nil {
		logs.ErrorLog.Printf("Invalid offset: %q", err)
		auth.JsRespond(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	profiles, err := modules.GetSearchprofile(query, offset)
	if err != nil {
		logs.ErrorLog.Printf("Error getting search profiles: %q", err)
		auth.JsRespond(w, "Failed to get search profiles", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profiles)
}
