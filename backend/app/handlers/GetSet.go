package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/logs"
	"social-network/app/modules"
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
	case "avatar":
		GetAvatarHandler(w, r, data.Sub)
	case "posts":
		GetPostsHandler(w, r, data.Sub)
	case "comments":
		GetCommentsHandler(w, r, data.Sub)
	case "groupMembers":
		GroupMembersHandler(w, r, data.Sub)
	case "groupFeeds":
		GroupFeedsHandler(w, r, data.Sub)
	case "groupImIn":
		GroupImInHandler(w, r, data.Sub)
	case "groupEvents":
		GroupEventsHandler(w, r, data.Sub)
	case "requests":
		GetRequestsHandler(w, r, data.Sub)
	case "userSeggestions":
		GetUserSuggestions(w, r, data.Sub)
	case "groupData":
		GetGroupDataHandler(w, r, data.Sub)
	case "search":
		GetSearchHandler(w, r, data.Sub)
	case "myFollowers":
		GetFollowersHandler(w, r, data.Sub)
	case "users":
		payload := r.Context().Value(auth.UserContextKey)
		data, ok := payload.(*jwt.JwtPayload)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		usernames, err := modules.GetUserNames(data.Sub)
		if err != nil {
			logs.ErrorLog.Printf("routes.go 60 %q", err.Error())
		}
		jsonData, _ := json.Marshal(usernames)
		w.Write(jsonData)
	case "dmhistory":
		target := r.URL.Query().Get("person_name")
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{error": "expected page value"}`)
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
		dms, err := modules.GetdmHistory(data.Sub, username, target, page)
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
		logs.ErrorLog.Println("Invalid request type")
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
	case "joinGroup":
		JoinGroup(w, r, data.Sub)
	case "eventCreation":
		GroupEventCreation(w, r, data.Sub)
	case "reactionEvents":
		GroupEventResponse(w, r, data.Sub)
	case "follow":
		FollowHandle(w, r, data.Sub)
	case "acceptFollow":
		FollowersAR(w, r, data.Sub)
	case "sendRequest":
		SendRequestHandler(w, r, data.Sub)
	case "like":
		LikeDislike(w, r, data.Sub)
	case "comment":
		CreateComment(w, r, data.Sub)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "page not found"}`)
		logs.ErrorLog.Printf("Invalid request to /set/: %s", r.PathValue("type"))
		return
	}
}
