package handlers

import (
	"encoding/json"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/modules"
	"social-network/app/structs"
)

func LikeDislike(w http.ResponseWriter, r *http.Request, uid int) {
	var LikeInfo structs.LikeInfo

	json.NewDecoder(r.Body).Decode(&LikeInfo)

	if !modules.LikeDislike(LikeInfo, uid) {
		auth.JsResponse(w, "Like/dislike failed", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Liked/disliked successfully",
	})
}
