package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
	"social-network/app/structs"
)

// function to create a post
func PostCreation(w http.ResponseWriter, r *http.Request, uid int) {
	var post structs.PostCreate

	json.NewDecoder(r.Body).Decode(&post)

	if !modules.InsertPost(post, uid, post.GroupId) {
		auth.JsResponse(w, "Post creation failed", http.StatusBadRequest)
		logs.ErrorLog.Println("Post creation failed")
		return
	}
	auth.JsResponse(w, "Post created successfully", http.StatusOK)
}

// function to get posts with filter of wich page to fetch
func GetPostsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	var dataToFetch map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataToFetch)
	if err != nil {
		logs.ErrorLog.Println("Error decoding request body:", err)
		auth.JsResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	forwhat, ok := dataToFetch["fetch"].(string)
	if !ok {
		logs.ErrorLog.Println("Invalid 'fetch' value:", dataToFetch["fetch"])
		auth.JsResponse(w, "Invalid 'fetch' value", http.StatusBadRequest)
		return
	}

	startFloat, ok := dataToFetch["start"].(float64)
	if !ok {
		logs.ErrorLog.Println("Invalid 'start' value:", dataToFetch["start"])
		auth.JsResponse(w, "Invalid 'start' value", http.StatusBadRequest)
		return
	}
	start := int(startFloat)

	userId := 0
	uidFloat, ok := dataToFetch["userId"].(float64)
	if ok {
		userId = int(uidFloat)
	}

	groupId := 0
	if gid, exists := dataToFetch["groupId"]; exists {
		if gidFloat, ok := gid.(float64); ok {
			groupId = int(gidFloat)
		}
	}

	switch forwhat {
	case "home":
		GetHomePostsHandler(w, r, uid, start)
	case "profile":
		GetProfilePostsHandler(w, r, uid, start, userId)
	case "group":
		GetGroupPostsHandler(w, r, uid, start, groupId)
	default:
		auth.JsResponse(w, "Invalid fetch type", http.StatusBadRequest)
		return
	}
}

// function to get posts of home page
func GetHomePostsHandler(w http.ResponseWriter, r *http.Request, uid, start int) {
	posts, err := modules.GetHomePosts(start, uid)
	if err != nil {
		auth.JsResponse(w, "Get Home Posts failed", http.StatusBadRequest)
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

// function to get posts of profile page
func GetProfilePostsHandler(w http.ResponseWriter, r *http.Request, uid, start, userId int) {
	var posts []structs.Post
	var err error
	if userId == 0 {
		posts, err = modules.GetOwnProfilePosts(start, uid)
		if err != nil {
			auth.JsResponse(w, "Get Own Profile Posts failed", http.StatusBadRequest)
			logs.ErrorLog.Println(err)
			return
		}
	} else {
		posts, err = modules.GetProfilePosts(start, uid, userId)
		if err != nil {
			auth.JsResponse(w, "Get Profile Posts failed", http.StatusBadRequest)
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

// function to get posts of group page
func GetGroupPostsHandler(w http.ResponseWriter, r *http.Request, uid, start, groupId int) {
	posts, err := modules.GetGroupPosts(start, uid, groupId)
	if err != nil {
		auth.JsResponse(w, "Get Group Posts failed", http.StatusBadRequest)
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

// function to get followers
func GetFollowersHandler(w http.ResponseWriter, r *http.Request, uid int) {
	rows, err := modules.DB.Query(`
		SELECT
		    p.id,
		    p.display_name
		FROM
		    follow f
		    JOIN profile p ON (
		        f.following_id = p.id
		        AND p.is_user = 1
		    )
		WHERE
		    f.follower_id = ?
	`, uid)
	if err != nil {
		logs.ErrorLog.Printf("GetFollowers query error: %q", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followers []structs.UsersGet
	for rows.Next() {
		var follower structs.UsersGet
		if err := rows.Scan(&follower.ID, &follower.Username); err != nil {
			logs.ErrorLog.Printf("Error scanning followers: %q", err.Error())
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		followers = append(followers, follower)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}
