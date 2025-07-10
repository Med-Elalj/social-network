package handlers

import (
	"encoding/json"
	"net/http"
	auth "social-network/app/Auth"
	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

func RequestsGet(w http.ResponseWriter, r *http.Request, uid int) {
	requests, err := modules.GetRequests(uid)
	if err != nil {
		logs.ErrorLog.Println("Error getting requests:", err)
		auth.JsRespond(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	response := make([]structs.RequestsGet, len(requests))
	for i, req := range requests {
		response[i] = structs.RequestsGet{
			ID:       req.ID,
			SenderId: req.SenderId,
			Towhat:   req.Towhat,
			Message:  req.Message,
			Username: req.Username,
			Avatar:   req.Avatar,
		}
	}

	json.NewEncoder(w).Encode(response)
}
