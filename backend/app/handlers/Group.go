package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
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
	err := modules.UpdatEventResp(event.ID, uid,event.Respond)
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
