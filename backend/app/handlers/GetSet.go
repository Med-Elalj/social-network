package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/modules"
	"social-network/server/logs"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(auth.UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Sorry something went wrong"}`)
		return
	}

	switch r.PathValue("type") {
	case "posts":
		GetPostsHandler(w, r, data.Sub)
	case "groupPosts":
		GetPostsHandler(w, r, data.Sub)
	case "groupMembers":
		GroupMembersHandler(w, r, data.Sub)
	case "groupFeeds":
		GroupFeedsHandler(w, r, data.Sub)
	case "groupToJoin":
		GroupToJoinHandler(w, r, data.Sub)
	case "groupImIn":
		GroupImInHandler(w, r, data.Sub)
	case "users":
		payload := r.Context().Value(auth.UserContextKey)
		data, ok := payload.(*jwt.JwtPayload)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		usernames, _ := modules.GetUserNames(data.Sub)
		jsonData, _ := json.Marshal(usernames)
		w.Write(jsonData)
	case "dmhistory":
		target := r.Header.Get("target")
		page, err := strconv.Atoi(r.Header.Get("page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{error": "expected page valur"}`)
			return
		}
		payload := r.Context().Value(auth.UserContextKey)
		data, ok := payload.(*jwt.JwtPayload)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Sorry something went wrong"}`)
			return
		}
		username := data.Username
		dms, err := modules.GetdmHistory(username, target, page)
		if err != nil {
			logs.ErrorLog.Printf("routes.go 69 %q", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Sorry something went wrong"}`)
			return
		}
		jsonData, _ := json.Marshal(dms)
		w.Write(jsonData)
	// 	// TODO get notifications
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
	case "follow":
		FollowersJoin(w, r, data.Sub)
	case "unfollow":
		FollowersLeave(w, r, data.Sub)
	case "acceptFollow":
		FollowersAccept(w, r, data.Sub)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "page not found"}`)
		logs.ErrorLog.Printf("Invalid request to /set/: %s", r.PathValue("type"))
		return
	}
}
