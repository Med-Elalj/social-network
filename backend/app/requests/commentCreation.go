package requests

// func CommentCreation(w http.ResponseWriter, r *http.Request, uid int) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		logs.Println("Error reading request body:", err)
// 		auth.JsRespond(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
// 		return
// 	}

// 	if len(body) == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintf(w, `{"error": "Request body cannot be empty"}`)
// 		return
// 	}

// 	var comment structs.CommentInfo

// 	structs.JsonRestrictedDecoder(body, &comment)
// 	if err := comment.Validate(); err != nil {
// 		logs.Println("Validation failed for comment content:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintf(w, `{"error": %q}`, err.Error())
// 		return
// 	}

// 	if !db.InsertComment(comment, uid) {
// 		structs.JsRespond(w, "Comment creation failed", http.StatusInternalServerError)
// 	}
// 	structs.JsRespond(w, "Comment posted successfully", http.StatusOK)
// }
