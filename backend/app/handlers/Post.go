package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/app/modules"
	"social-network/app/structs"
	"social-network/server/logs"
)

func PostCreation(w http.ResponseWriter, r *http.Request, uid int) {
	var post structs.PostCreate

	json.NewDecoder(r.Body).Decode(&post)

	if !modules.InsertPost(post, uid, nil) {
		structs.JsRespond(w, "Post creation failed", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post Added successfully",
	})
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var dataToFetch map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataToFetch)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse start
	startFloat, ok := dataToFetch["start"].(float64)
	if !ok {
		http.Error(w, "'start' must be a number", http.StatusBadRequest)
		return
	}
	start := int(startFloat)

	// Parse groupId, default to 0 if not present or invalid
	groupId := 0
	if gid, exists := dataToFetch["groupId"]; exists {
		if gidFloat, ok := gid.(float64); ok {
			groupId = int(gidFloat)
		}
	}

	posts, err := modules.GetPosts(start, uid, groupId)
	if err != nil {
		structs.JsRespond(w, "Get Posts failed", http.StatusBadRequest)
		logs.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Posts fetched successfully",
		"posts":   posts,
	})
}