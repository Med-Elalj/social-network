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
	var bodyRequest struct {
		Type      int  `json:"type"`
		IsSpecial bool `json:"is_special"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bodyRequest); err != nil {
		logs.ErrorLog.Printf("Failed to decode request body: %q", err)
		auth.JsResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	requests, err := modules.GetRequests(uid, bodyRequest.Type, bodyRequest.IsSpecial)
	if err != nil {
		logs.ErrorLog.Println("Error getting requests:", err)
		auth.JsResponse(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

func SendRequestHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var bodyRequest struct {
		Target     int  `json:"target"`
		Type       int  `json:"type"`
		ReceiverId int  `json:"receiver_id"`
		Is_public  bool `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bodyRequest); err != nil {
		logs.ErrorLog.Printf("Failed to decode request body: %q", err)
		auth.JsResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("request body:", bodyRequest)

	if bodyRequest.ReceiverId == 0 {
		if bodyRequest.Type == 0 && bodyRequest.Is_public {
			err := modules.InsertFollow(uid, bodyRequest.Target)
			if err != nil {
				logs.ErrorLog.Println("Error inserting follow relationship:", err)
				auth.JsResponse(w, "Failed to send follow request", http.StatusInternalServerError)
				return
			}
			auth.JsResponse(w, "Follow request sent successfully", http.StatusOK)
			return
		}

		if err := modules.InsertRequest(uid, bodyRequest.ReceiverId, bodyRequest.Target, bodyRequest.Type); err != nil {
			logs.ErrorLog.Println("Error sending request:", err)
			auth.JsResponse(w, "Failed to send request", http.StatusInternalServerError)
			return
		}
	} else {
		if err := modules.InsertGroupRequestFromUser(uid, bodyRequest.Target, bodyRequest.ReceiverId); err != nil {
			logs.ErrorLog.Println("Error sending request:", err)
			auth.JsResponse(w, "Failed to send request", http.StatusInternalServerError)
			return
		}
	}

	auth.JsResponse(w, "Request sent successfully", http.StatusOK)
}
