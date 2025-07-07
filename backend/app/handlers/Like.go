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

	if !modules.LikeDeslike(LikeInfo, uid) {
		auth.JsRespond(w, "Like/deslike failed", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Liked/desliked successfully",
	})
}
