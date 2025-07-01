package handlers

import (
	"net/http"
)

type Profile struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Description string `json:"description,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Followers   int    `json:"followers"`
	Followed    int    `json:"followed"`
	Online      bool   `json:"online"`
}

func ProfileHandler(w http.ResponseWriter, r *http.Request, uid int) {
	// get profile information from database
	// respond with profile data in JSON format
}
