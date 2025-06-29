package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/app/modules"
	"social-network/app/structs"
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
