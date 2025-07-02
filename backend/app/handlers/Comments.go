package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/modules"
	"social-network/app/structs"
)

func CreateComment(w http.ResponseWriter, r *http.Request,uid int) {
	var comment structs.CommentInfo

	json.NewDecoder(r.Body).Decode(&comment)

	if !modules.InsertComment(comment, uid) {
		auth.JsRespond(w, "Comment creation failed", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Comment Added successfully",
	})
}
