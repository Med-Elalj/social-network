package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
	"social-network/app/structs"
)


func getfollowers(w http.ResponseWriter, r *http.Request, uid int) {
	var dataToFetch map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataToFetch)
	if err != nil {
		logs.ErrorLog.Println("Error decoding request body:", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startFloat, ok := dataToFetch["start"].(float64)
	if !ok {
		logs.ErrorLog.Println("Invalid 'start' value:", dataToFetch["start"])
		auth.JsRespond(w, "Invalid 'start' value", http.StatusBadRequest)
		return
	}
	start := int(startFloat)

	followers, err := modules.GetFollowers(start,uid)
	if err != nil {
		auth.JsRespond(w, "Get Followers failed", http.StatusBadRequest)
		logs.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Followers fetched successfully",
		"followers": followers,
	})
}

func PostCreation(w http.ResponseWriter, r *http.Request, uid int) {
	var post structs.PostCreate

	json.NewDecoder(r.Body).Decode(&post)

	// fmt.Println(post.GroupId)
	// fmt.Println(uid)
	// fmt.Println(post.Privacy)
	// fmt.Println(post.Content)

	if !modules.InsertPost(post, uid, post.GroupId) {
		auth.JsRespond(w, "Post creation failed", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post Added successfully",
	})
}

func GetGroupPostsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var dataToFetch map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataToFetch)
	if err != nil {
		logs.ErrorLog.Println("Error decoding request body:", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	startFloat, ok := dataToFetch["start"].(float64)
	if !ok {
		logs.ErrorLog.Println("Invalid 'start' value:", dataToFetch["start"])
		auth.JsRespond(w, "Invalid 'start' value", http.StatusBadRequest)
		return
	}
	start := int(startFloat)

	groupId := 0
	if gid, exists := dataToFetch["groupId"]; exists {
		if gidFloat, ok := gid.(float64); ok {
			groupId = int(gidFloat)
		}
	}

	posts, err := modules.GetGroupPosts(start, uid, groupId)
	if err != nil {
		auth.JsRespond(w, "Get Group Posts failed", http.StatusBadRequest)
		logs.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Group Posts fetched successfully",
		"posts":   posts,
	})
}

func GetProfilePostsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var dataToFetch map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataToFetch)
	if err != nil {
		logs.ErrorLog.Println("Error decoding request body:", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startFloat, ok := dataToFetch["start"].(float64)
	if !ok {
		logs.ErrorLog.Println("Invalid 'start' value:", dataToFetch["start"])
		auth.JsRespond(w, "Invalid 'start' value", http.StatusBadRequest)
		return
	}
	start := int(startFloat)

	userId := 0
	if userIdMap, ok := dataToFetch["userId"].(map[string]interface{}); ok {
		if uidFloat, ok := userIdMap["userId"].(float64); ok {
			userId = int(uidFloat)
		}
	}

	var posts []structs.Post
	// var err error
	if userId == 0 {
		posts, err = modules.GetOwnProfilePosts(start, uid)
		if err != nil {
			auth.JsRespond(w, "Get Own Profile Posts failed", http.StatusBadRequest)
			logs.ErrorLog.Println(err)
			return
		}
	} else {
		posts, err = modules.GetProfilePosts(start, uid, userId)
		if err != nil {
			auth.JsRespond(w, "Get Profile Posts failed", http.StatusBadRequest)
			logs.ErrorLog.Println(err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile Posts fetched successfully",
		"posts":   posts,
	})
}

func GetHomePostsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var dataToFetch map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataToFetch)
	if err != nil {
		logs.ErrorLog.Println("Error decoding request body:", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startFloat, ok := dataToFetch["start"].(float64)
	if !ok {
		logs.ErrorLog.Println("Invalid 'start' value:", dataToFetch["start"])
		auth.JsRespond(w, "Invalid 'start' value", http.StatusBadRequest)
		return
	}
	start := int(startFloat)

	posts, err := modules.GetHomePosts(start, uid)
	if err != nil {
		auth.JsRespond(w, "Get Home Posts failed", http.StatusBadRequest)
		logs.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Home Posts fetched successfully",
		"posts":   posts,
	})
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var dataToFetch map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataToFetch)
	if err != nil {
		logs.ErrorLog.Println("Error decoding request body:", err)
		auth.JsRespond(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	forwhat, ok := dataToFetch["fetch"].(string)
	if !ok {
		logs.ErrorLog.Println("Invalid 'fetch' value:", dataToFetch["fetch"])
		auth.JsRespond(w, "Invalid 'fetch' value", http.StatusBadRequest)
		return
	}
	switch forwhat {
	case "home":
		GetHomePostsHandler(w, r, uid)
	case "profile":
		GetProfilePostsHandler(w, r, uid)
	case "group":
		GetGroupPostsHandler(w, r, uid)
	default:
		auth.JsRespond(w, "Invalid fetch type", http.StatusBadRequest)
		return
	}
}
