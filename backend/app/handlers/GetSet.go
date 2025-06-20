package handlers

import (
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// offset := r.URL.Query().Get("offset")

	switch r.PathValue("type") {
	case "comments":
		// pid := r.URL.Query().Get("pid")
		// comments, _ := db.GetComments(pid)
		// jsoncomment, _ := json.Marshal(comments)
		// w.Write([]byte(jsoncomment))
	case "posts":
		// posts, _ := db.GetPosts(offset)
		// jsonData, _ := json.Marshal(posts)
		// w.Write(jsonData)
	case "users":
		// payload := r.Context().Value(security.UserContextKey)
		// data, ok := payload.(*jwt.JwtPayload)
		// if !ok {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprintf(w, `{"error": "Sorry something went wrong"}`)
		// 	return
		// }
		// usernames, _ := db.GetUserNames(data.Sub)
		// jsonData, _ := json.Marshal(usernames)
		// w.Write(jsonData)
	case "dmhistory":
		// target := r.Header.Get("target")
		// page := r.Header.Get("page")
		// payload := r.Context().Value(security.UserContextKey)
		// data, ok := payload.(*jwt.JwtPayload)
		// if !ok {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprintf(w, `{"error": "Sorry something went wrong"}`)
		// 	return
		// }
		// username := data.Username
		// dms, err := db.GetdmHistory(username, target, page)
		// if err != nil {
		// 	logs.Errorf("routes.go 69 %q", err.Error())
		// }
		// jsonData, _ := json.Marshal(dms)
		// w.Write(jsonData)
		// TODO get notifications
	default:
		// w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(w, `{"error": "Invalid request type"}`)
		// return
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// payload := r.Context().Value(security.UserContextKey)
	// data, ok := payload.(*jwt.JwtPayload)
	// if !ok {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, `{"error": "Sorry something went wrong"}`)
	// 	return
	// }
	switch r.PathValue("type") {
	case "Post":
		// requests.PostCreation(w, r, data.Sub)
	case "Comment":
		// requests.CommentCreation(w, r, data.Sub)
		// TODO: set follow, profile
		// case "Follow":
		// 	requests.FollowCreation(w, r, data.Sub)
		// case "Profile":
		// 	requests.ProfileUpdate(w, r, data.Sub)
	case "Upload":
		// upload.UploadHandler(w, r, data.Sub)
	default:
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(w, `{"error": "page not found"}`)
		// logs.Errorf("Invalid request to /set/: %s", r.PathValue("type"))
		// return
	}
}
