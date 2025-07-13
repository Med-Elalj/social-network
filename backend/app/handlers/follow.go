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

	var bodyRequest BodyRequest

	err := json.NewDecoder(r.Body).Decode(&bodyRequest)

	if err != nil {
		logs.ErrorLog.Println("invalid request body:", err)
		auth.JsRespond(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := modules.UserFollow(uid, bodyRequest.Target, bodyRequest.Status); err != nil {
		logs.ErrorLog.Println("Error inserting follow relationship:", err)
		auth.JsRespond(w, "group joining failed", http.StatusInternalServerError)
	}

	// TODO notif to group creator
	auth.JsRespond(w, "user send req group successfully", http.StatusOK)
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
		Id     int    `json:"sender"`
		Target int    `json:"target"`
		Status string `json:"status"`
	}

	var bodyRequest BodyRequest

	err := json.NewDecoder(r.Body).Decode(&bodyRequest)
	if err != nil {
		logs.ErrorLog.Println("invalid request id:", err)
		auth.JsRespond(w, "invalid request id", http.StatusBadRequest)
		return
	}

	if bodyRequest.Status == "accept" {
		if err := modules.InsertFollow(bodyRequest.Id, bodyRequest.Target); err != nil {
			logs.ErrorLog.Println("Error accepting follow relationship:", err)
			auth.JsRespond(w, "group accepting failed", http.StatusInternalServerError)
			return
		}
	}

	if err := modules.DeleteRequest(bodyRequest.Id, uid, bodyRequest.Target); err != nil {
		logs.ErrorLog.Println("Error rejecting follow relationship:", err)
		auth.JsRespond(w, "group rejecting failed", http.StatusInternalServerError)
	}

	auth.JsRespond(w, fmt.Sprintf("user %sed request", bodyRequest.Status), http.StatusOK)

	// TODO notif to group creator

}
