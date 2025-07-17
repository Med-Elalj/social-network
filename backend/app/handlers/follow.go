package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
	"social-network/app/structs"
)

// needs header "follow_target" the id of the profile you want to follow
func FollowHandle(w http.ResponseWriter, r *http.Request, uid int) {
	type BodyRequest struct {
		Target int    `json:"target"`
		Status string `json:"status"`
	}
	type ResponseBody struct {
		NewStatus string `json:"new_status"`
	}

	var bodyRequest BodyRequest
	var responseBody ResponseBody

	err := json.NewDecoder(r.Body).Decode(&bodyRequest)
	if err != nil {
		logs.ErrorLog.Println("invalid request body:", err)
		auth.JsRespond(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if responseBody.NewStatus, err = modules.UserFollow(uid, bodyRequest.Target, bodyRequest.Status); err != nil {
		logs.ErrorLog.Println("Error inserting follow relationship:", err)
		auth.JsRespond(w, "error sent request, try again on another time", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseBody)

	// TODO notif to group creator
}

// needs header "follow_target" the id of the profile you want to unfollow
func FollowersLeave(w http.ResponseWriter, r *http.Request, uid int) {
	gid, err := strconv.Atoi(r.Header.Get("follow_target"))
	if err != nil {
		logs.ErrorLog.Println("Error converting group ID:", err)
		auth.JsRespond(w, "group id is required", http.StatusBadRequest)
		return
	}

	if err := modules.DeleteFollow(uid, gid); err != nil {
		logs.ErrorLog.Println("Error deleting follow relationship:", err)
		auth.JsRespond(w, "group leaving failed", http.StatusInternalServerError)
		return
	}

	// TODO notif to group creator
	auth.JsRespond(w, "user left group successfully", http.StatusOK)
}

// needs header "follower_target" the id of the follower you want to accept
// needs header "group_target" the id of the group you want to accept | 0 or unset means accepting to personal profile
func FollowersAccept(w http.ResponseWriter, r *http.Request, uid int) {
	gid, err := strconv.Atoi(r.Header.Get("group_target"))
	if err != nil || gid < 0 {
		gid = 0
	}
	folower_id, err := strconv.Atoi(r.Header.Get("follower_target"))
	if err != nil {
		logs.ErrorLog.Println("Error converting group ID:", err)
		auth.JsRespond(w, "follower_id as header is required", http.StatusBadRequest)
		return
	}
	if err := modules.AcceptFollow(uid, gid, folower_id); err != nil {
		logs.ErrorLog.Println("Error accepting follow relationship:", err)
		auth.JsRespond(w, "group accepting failed", http.StatusInternalServerError)
		return
	}

	// TODO notif to group creator
	auth.JsRespond(w, "user accepted group successfully", http.StatusOK)
}

func GetFollowRequests(w http.ResponseWriter, r *http.Request, uid int) {
	users, err := modules.GetFollowRequests(uid)
	if err != nil {
		logs.ErrorLog.Println("Error get follow requests from Db:", err)
		auth.JsRespond(w, "get follow requests failed", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string][]structs.Gusers{
		"users": users,
	})
}

func FollowersAR(w http.ResponseWriter, r *http.Request, uid int) {
	type BodyRequest struct {
		Id        int    `json:"sender"`
		Target    int    `json:"target"`
		Status    string `json:"status"`
		Type      int    `json:"type"`
		IsSpecial bool   `json:"isSpecial"`
	}
	type ResponseBody struct {
		NewStatus string `json:"new_status"`
		Message   string `json:"message"`
	}

	var err error
	var targetIsPublic bool
	var bodyRequest BodyRequest
	var responseBody ResponseBody

	err = json.NewDecoder(r.Body).Decode(&bodyRequest)
	if err != nil {
		logs.ErrorLog.Println("invalid request id:", err)
		auth.JsRespond(w, "invalid request id", http.StatusBadRequest)
		return
	}
	// Get target user's privacy status before processing
	if bodyRequest.Type == 0 {
		err = modules.DB.QueryRow(`SELECT is_public FROM profile WHERE id = ?`, bodyRequest.Id).Scan(&targetIsPublic)
		if err != nil {
			logs.ErrorLog.Println("Error getting target user privacy:", err)
			auth.JsRespond(w, "error processing request", http.StatusInternalServerError)
			return
		}
	}

	if bodyRequest.Status == "accept" {
		if bodyRequest.Type == 0 {
			bodyRequest.Target = uid
		}

		// First, insert the follow relationship (sender follows target)
		if bodyRequest.IsSpecial {
			err = modules.InsertFollow(uid, bodyRequest.Target)
		} else {
			err = modules.InsertFollow(bodyRequest.Id, bodyRequest.Target)
		}

		if err != nil {
			logs.ErrorLog.Println("Error accepting follow relationship:", err)
			auth.JsRespond(w, "follow accepting failed", http.StatusInternalServerError)
			return
		}

		// If it's a follow type request and the requester is public,
		// we should also make the accepter follow back
		if bodyRequest.Type == 0 && targetIsPublic {
			// Check if the accepter already follows the requester
			var alreadyFollowing bool
			err = modules.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM follow WHERE follower_id = ? AND following_id = ?)`,
				uid, bodyRequest.Id).Scan(&alreadyFollowing)
			if err != nil {
				logs.ErrorLog.Println("Error checking existing follow:", err)
				// Continue without follow back, at least the main follow worked
			} else if !alreadyFollowing {
				// Make the accepter follow back
				if err := modules.InsertFollow(uid, bodyRequest.Id); err != nil {
					logs.ErrorLog.Println("Error inserting follow back:", err)
					// Continue, the main follow relationship was created
				}
			}
		}
	}

	err = modules.DeleteRequest(bodyRequest.Id, uid, bodyRequest.Target, bodyRequest.Type)
	if err != nil {
		logs.ErrorLog.Println("Error deleting follow request:", err)
		auth.JsRespond(w, "error processing request", http.StatusInternalServerError)
		return
	}

	// Only get relationship status for follow type requests
	if bodyRequest.Type == 0 {
		responseBody.NewStatus, err = modules.GetRelationship(uid, bodyRequest.Id)
		if err != nil {
			logs.ErrorLog.Println("Error getting relationship:", err)
			// Don't fail the request, just set a default status
			responseBody.NewStatus = "follow"
		}

	}

	responseBody.Message = fmt.Sprintf("Request %sed successfully", bodyRequest.Status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseBody)
}

func GetUserSuggestions(w http.ResponseWriter, r *http.Request, uid int) {
	is_user, err := strconv.Atoi(r.URL.Query().Get("is_user"))
	if err != nil && is_user != 0 && is_user != 1 {
		logs.ErrorLog.Println("invalid is_user request: ", is_user, err)
		auth.JsRespond(w, "invalid is_user request", http.StatusBadRequest)
		return
	}
	users, err := modules.GetSuggestions(uid, is_user)
	if err != nil {
		logs.ErrorLog.Println("Error getting segguestions: ", err)
		auth.JsRespond(w, "error getting segguestions: ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(users)
}
