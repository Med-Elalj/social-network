package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"social-network/app/modules"
	"social-network/app/security"
	"social-network/app/security/jwt"
	"social-network/app/structs"
	"social-network/server/logs"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var group structs.PostGet
	json.NewDecoder(r.Body).Decode(&group) // well be flexible to take any data

	switch r.PathValue("type") {
	case "posts":
		posts, _ := modules.GetPosts(int(group.Start), int(group.UserId), int(group.GroupId))

		json.NewEncoder(w).Encode(map[string][]structs.Post{
			"posts": posts,
		})
	case "groupPosts":
		posts, _ := modules.GetPosts(int(group.Start), int(group.UserId), int(group.GroupId))

		json.NewEncoder(w).Encode(map[string][]structs.Post{
			"posts": posts,
		})
	case "groupMembers":
		groupid, err := strconv.Atoi(r.URL.Query().Get("gid"))
		if err != nil {
			logs.Fatalln(err.Error())
		}
		members, _ := modules.GetMembers(groupid)
		json.NewEncoder(w).Encode(map[string][]structs.Gusers{
			"members": members,
		})
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid request type"}`)
		return
	}
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(security.UserContextKey)
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
