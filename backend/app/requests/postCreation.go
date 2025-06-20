package requests

// func PostCreation(w http.ResponseWriter, r *http.Request, uid int) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		logs.Println("Error reading request body:", err)
// 		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
// 		return
// 	}

// 	if len(body) == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintf(w, `{"error": "Request body cannot be empty"}`)
// 		return
// 	}

// 	var post structs.PostCreate

// 	structs.JsonRestrictedDecoder(body, &post)

// 	post.ParseCategories()
// 	if err := post.Validate(); err != nil {
// 		logs.Println("Validation failed for post title:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintf(w, `{"error": %q}`, err.Error())
// 		return
// 	}
// 	if !db.InsertPost(post, uid) {
// 		structs.JsRespond(w, "Post creation failed", http.StatusBadRequest)
// 	}
// 	structs.JsRespond(w, "Post created successfully", http.StatusOK)
// }
