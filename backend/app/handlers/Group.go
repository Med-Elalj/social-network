package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
	"social-network/app/structs"
)

func GroupEventsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var groupId int
	json.NewDecoder(r.Body).Decode(&groupId)

	events, err := modules.GetEvents(groupId, uid)
	if err != nil {
		auth.JsRespond(w, "Failed to get group events", http.StatusBadRequest)
		logs.ErrorLog.Println("Error getting group events:", err)
		return
	}
	json.NewEncoder(w).Encode(map[string][]structs.GroupEvent{
		"events": events,
	})
}

func UpdateResponseHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var event structs.GroupEvent
	json.NewDecoder(r.Body).Decode(&event)
	err := modules.UpdatEventResp(event.ID, uid, event.Respond)
	if err != nil {
		auth.JsRespond(w, "Failed to update response", http.StatusBadRequest)
		logs.ErrorLog.Println("Error updating response:", err)
		return
	}
}

func GroupEventCreation(w http.ResponseWriter, r *http.Request, uid int) {
	var event structs.GroupEvent

	json.NewDecoder(r.Body).Decode(&event)
	lastID, err := modules.Insertevent(event, uid)
	if err != nil {
		auth.JsRespond(w, "Failed to create event", http.StatusBadRequest)
		logs.ErrorLog.Println("Error inserting event into database:", err)
		return
	}

	err = modules.InsertUserEvent(lastID, uid, true)
	if err != nil {
		auth.JsRespond(w, "Failed to add creator to the event", http.StatusBadRequest)
		logs.ErrorLog.Println("Error inserting event into database:", err)
		return
	}

	auth.JsRespond(w, "event adding successfully", http.StatusOK)
}

func GroupCreation(w http.ResponseWriter, r *http.Request, uid int) {
	var group structs.Group

	json.NewDecoder(r.Body).Decode(&group)
	err := modules.InsertGroup(group, uid)
	if err != nil {
		auth.JsRespond(w, "group creation failed", http.StatusInternalServerError)
		logs.ErrorLog.Println("Error inserting group into database:", err)
		return
	}
	auth.JsRespond(w, "group Created successfully", http.StatusOK)
}

func GroupToJoinHandler(w http.ResponseWriter, r *http.Request, uid int) {
	groups, err := modules.GetGroupToJoin(uid)
	if err != nil {
		auth.JsRespond(w, "Failed to get groups to", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string][]structs.GroupGet{
		"groups": groups,
	})
}

func GroupMembersHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var groupId int

	json.NewDecoder(r.Body).Decode(&groupId)

	members, err := modules.GetMembers(groupId)
	if err != nil {
		auth.JsRespond(w, "Failed to get group members", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string][]structs.Gusers{
		"members": members,
	})
}

func GroupFeedsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	posts, err := modules.GetGroupFeed(uid)
	if err != nil {
		auth.JsRespond(w, "Failed to get group feeds", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string][]structs.Post{
		"posts": posts,
	})
}

func GroupImInHandler(w http.ResponseWriter, r *http.Request, uid int) {
	groups, err := modules.GetGroupImIn(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get groups"})
	}
	json.NewEncoder(w).Encode(map[string][]structs.GroupGet{
		"groups": groups,
	})
}

func GetGroupDataHandler(w http.ResponseWriter, r *http.Request, uid int) {
	type groupRequest struct {
		GroupName string `json:"groupName"`
	}

	var req groupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	groupName := req.GroupName

	var groupData structs.GroupGet

	query := `SELECT id, display_name, avatar, description, is_Public
	          FROM profile
	          WHERE display_name = ? AND is_user = 0;`
	err := modules.DB.QueryRow(query, groupName).Scan(&groupData.ID, &groupData.GroupName, &groupData.Avatar, &groupData.Description, &groupData.Privacy)
	if err != nil {
		logs.ErrorLog.Printf("Group Data Handler scan error: %q", err.Error())
		auth.JsRespond(w, "Failed to get group data", http.StatusBadRequest)
		return
	}

	memberQuery := `SELECT COUNT(*) FROM follow WHERE following_id = ?`
	var memberCount int
	err = modules.DB.QueryRow(memberQuery, groupData.ID).Scan(&memberCount)
	if err != nil {
		logs.ErrorLog.Printf("Group Data Handler scan error: %q", err.Error())
		auth.JsRespond(w, "Failed to get group member count", http.StatusBadRequest)
		return
	}

	groupData.MemberCount = memberCount + 1

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupData)
}

func JoinGroup(w http.ResponseWriter, r *http.Request, uid int) {
	var bodyRequest struct {
		GroupId int `json:"groupId"`
	}

	err := json.NewDecoder(r.Body).Decode(&bodyRequest)
	if err != nil {
		logs.ErrorLog.Printf("invalid request body: %q", err)
		auth.JsRespond(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = modules.InsertRequest(uid, bodyRequest.GroupId, 1)
	if err != nil {
		auth.JsRespond(w, "error inserting new request", http.StatusInternalServerError)
		return
	}

	auth.JsRespond(w, "join request sented succeffully", http.StatusOK)
}
