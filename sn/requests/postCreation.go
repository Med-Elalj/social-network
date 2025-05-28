package requests

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"

	"social-network/sn/db"
	"social-network/sn/structs"
)

func PostCreation(w http.ResponseWriter, r *http.Request, uid int) {
	var post structs.PostInfo
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		structs.JsRespond(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	categories := strings.Split(strings.Join(strings.Fields(html.EscapeString(post.Category)), " "), ",")
	if len(categories) > 3 || len(categories) < 1 {
		structs.JsRespond(w, "Invalid category selection", http.StatusBadRequest)
		return
	}

	// Validate basic input
	if len(post.Title) < 3 || len(post.Content) < 10 {
		structs.JsRespond(w, "Title and content required", http.StatusBadRequest)
		return
	}
	if len(post.Content) > 1500 || len(post.Title) > 30 {
		structs.JsRespond(w, "Title or content too long", http.StatusBadRequest)
		return
	}
	if !db.InsertPost(post, categories, uid) {
		structs.JsRespond(w, "Post creation failed", http.StatusBadRequest)
	}
	structs.JsRespond(w, "Post created successfully", http.StatusOK)
}
