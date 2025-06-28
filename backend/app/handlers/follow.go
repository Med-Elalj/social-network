package handlers

import (
	"net/http"
	"strconv"

	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

// needs header "follow_target" the id of the profile you want to follow
func FollowersJoin(w http.ResponseWriter, r *http.Request, uid int) {
	gid, err := strconv.Atoi(r.Header.Get("follow_target"))
	if err != nil {
		logs.Println("Error converting group ID:", err)
		structs.JsRespond(w, "group id is required", http.StatusBadRequest)
		return
	}

	if err := modules.InsertFollow(uid, gid); err != nil {
		logs.Println("Error inserting follow relationship:", err)
		structs.JsRespond(w, "group joining failed", http.StatusInternalServerError)
	}

	// TODO notif to group creator
	structs.JsRespond(w, "user send req group successfully", http.StatusOK)
}

// needs header "follow_target" the id of the profile you want to unfollow
func FollowersLeave(w http.ResponseWriter, r *http.Request, uid int) {
	gid, err := strconv.Atoi(r.Header.Get("follow_target"))
	if err != nil {
		logs.Println("Error converting group ID:", err)
		structs.JsRespond(w, "group id is required", http.StatusBadRequest)
		return
	}

	if err := modules.DeleteFollow(uid, gid); err != nil {
		logs.Println("Error deleting follow relationship:", err)
		structs.JsRespond(w, "group leaving failed", http.StatusInternalServerError)
		return
	}

	// TODO notif to group creator
	structs.JsRespond(w, "user left group successfully", http.StatusOK)
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
		logs.Println("Error converting group ID:", err)
		structs.JsRespond(w, "follower_id as header is required", http.StatusBadRequest)
		return
	}
	if err := modules.AcceptFollow(uid, gid, folower_id); err != nil {
		logs.Println("Error accepting follow relationship:", err)
		structs.JsRespond(w, "group accepting failed", http.StatusInternalServerError)
		return
	}

	// TODO notif to group creator
	structs.JsRespond(w, "user accepted group successfully", http.StatusOK)
}
