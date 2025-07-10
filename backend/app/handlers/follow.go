package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	auth "social-network/app/Auth"
	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

// needs header "follow_target" the id of the profile you want to follow
func FollowersJoin(w http.ResponseWriter, r *http.Request, uid int) {
	gid, err := strconv.Atoi(r.Header.Get("follow_target"))
	if err != nil {
		logs.ErrorLog.Println("Error converting group ID:", err)
		auth.JsRespond(w, "group id is required", http.StatusBadRequest)
		return
	}

	if err := modules.InsertFollow(uid, gid); err != nil {
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

func followersList(w http.ResponseWriter, r *http.Request, uid int) {
	w.Header().Set("Content-Type", "application/json")
	followers, err := modules.getfollowers(uid)
	if err != nil {
		logs.ErrorLog.Println("Error getting followers:", err)
		auth.JsRespond(w, "failed to get followers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string][]structs.UsersGet{
		"followers": followers,
	})
}

func followingList(w http.ResponseWriter, r *http.Request, uid int) {
	w.Header().Set("Content-Type", "application/json")
	following, err := modules.getfollowing(uid)
	if err != nil {
		logs.ErrorLog.Println("Error getting following:", err)
		auth.JsRespond(w, "failed to get following", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string][]structs.UsersGet{
		"following": following,
	})
}


func followcreateHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var follow structs.FollowReq
	if err := json.NewDecoder(r.Body).Decode(&follow); err != nil {
		logs.ErrorLog.Println("Error decoding follow request:", err)
		auth.JsRespond(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if follow.FollowingId == 0 {
		return
	}

	if err := modules.userfollow(uid, follow.FollowingId); err != nil {
		return
	}

	auth.JsRespond(w, "follow created successfully", http.StatusOK)
}