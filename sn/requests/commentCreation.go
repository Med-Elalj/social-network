package requests

import (
	"encoding/json"
	"net/http"

	"social-network/sn/db"
	"social-network/sn/structs"
)

func CommentCreation(w http.ResponseWriter, r *http.Request, uid int) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var comment structs.CommentInfo
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(comment.Content) > 1000 {
		structs.JsRespond(w, "Comment length exceeds the allowed ", http.StatusBadRequest)
		return
	}
	if !db.InsertComment(comment, uid) {
		structs.JsRespond(w, "Comment creation failed", http.StatusInternalServerError)
	}
	structs.JsRespond(w, "Comment posted successfully", http.StatusOK)
}
