package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
	"social-network/app/structs"
)

func GroupRequestsHandler(w http.ResponseWriter, r *http.Request, uid int) {
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

	response := make([]structs.RequestsGet, len(requests))
	for i, req := range requests {
		response[i] = structs.RequestsGet{
			SenderId:    req.SenderId,
			GroupId:     req.GroupId,
			GroupName:   req.GroupName,
			GroupAvatar: req.GroupAvatar,
			Type:        req.Type,
			Message:     req.Message,
			Username:    req.Username,
			Avatar:      req.Avatar,
		}
	}

	json.NewEncoder(w).Encode(response)
}
