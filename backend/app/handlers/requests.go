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
		Type int `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bodyRequest); err != nil {
		fmt.Println(err)
		logs.ErrorLog.Printf("Failed to decode request body: %q", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest) // error
		return
	}

	requests, err := modules.GetRequests(uid, bodyRequest.Type)
	if err != nil {
		logs.ErrorLog.Println("Error getting requests:", err)
		auth.JsRespond(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	fmt.Println(requests)

	json.NewEncoder(w).Encode(requests)
}

func SendRequestHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var bodyRequest struct {
		Target    int  `json:"target"`
		Type      int  `json:"type"`
		Is_public bool `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bodyRequest); err != nil {
		logs.ErrorLog.Printf("Failed to decode request body: %q", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if bodyRequest.Type == 0 && bodyRequest.Is_public {
		err := modules.InsertFollow(uid, bodyRequest.Target)
		if err != nil {
			logs.ErrorLog.Println("Error inserting follow relationship:", err)
			auth.JsRespond(w, "Failed to send follow request", http.StatusInternalServerError)
			return
		}
		auth.JsRespond(w, "Follow request sent successfully", http.StatusOK)
		return
	}

	if err := modules.InsertRequest(uid, bodyRequest.Target, bodyRequest.Target, bodyRequest.Type); err != nil {
		logs.ErrorLog.Println("Error sending request:", err)
		auth.JsRespond(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	auth.JsRespond(w, "Request sent successfully", http.StatusOK)
}
