package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dataToFetch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&dataToFetch); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	switch r.PathValue("type") {
	case "posts":
		start := int(dataToFetch["start"].(float64))
		userId := int(dataToFetch["userId"].(float64))
		groupId := int(dataToFetch["groupId"].(float64))

		posts, err := modules.GetPosts(start, userId, groupId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get posts"})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"posts": posts,
		})
	case "groupPosts":
		start := int(dataToFetch["start"].(float64))
		userId := int(dataToFetch["userId"].(float64))
		groupId := int(dataToFetch["groupId"].(float64))

		posts, err := modules.GetPosts(int(start), int(userId), int(groupId))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get groups"})
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string][]structs.Post{
			"posts": posts,
		})
	case "groupMembers":
		groupId := int(dataToFetch["groupId"].(float64))

		members, err := modules.GetMembers(groupId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get groups"})
		}
		json.NewEncoder(w).Encode(map[string][]structs.Gusers{
			"members": members,
		})
	case "groupFeeds":
		userId := int(dataToFetch["userId"].(float64))
		posts, err := modules.GetGroupFeed(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get groups"})
		}

		json.NewEncoder(w).Encode(map[string][]structs.Post{
			"posts": posts,
		})
	case "groupToJoin":
		userId := int(dataToFetch["userId"].(float64))

		groups, err := modules.GetGroupToJoin(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get groups"})
			return
		}
		json.NewEncoder(w).Encode(map[string][]structs.GroupGet{
			"groups": groups,
		})
	case "groupImIn":
		userId := int(dataToFetch["userId"].(float64))

		groups, err := modules.GetGroupImIn(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get groups"})
		}
		json.NewEncoder(w).Encode(map[string][]structs.GroupGet{
			"groups": groups,
		})
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request type"})
	}
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(auth.UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Sorry something went wrong"}`)
		return
	}
	switch r.PathValue("type") {
	case "Post":
		PostCreation(w, r, data.Sub)
	case "GroupCreation":
		GroupCreation(w, r, data.Sub)
	case "GroupJoin":
		GroupJoin(w, r, data.Sub)
	case "GroupUactive":
		GroupUactive(w, r, data.Sub)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "page not found"}`)
		logs.Errorf("Invalid request to /set/: %s", r.PathValue("type"))
		return
	}
}
