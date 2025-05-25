package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"social-network/sn/db"
	"social-network/sn/hub"
	"social-network/sn/structs"
)

func getPosts() (*sql.Rows, error) {
	rows, err := db.DB.Query(`
    SELECT 
        posts.id,
        posts.title, 
        posts.created_at, 
        users.username, 
        GROUP_CONCAT(category.name, ', ') AS categories, 
			COALESCE(comment.comment_count, 0) AS comment,
    FROM posts
    INNER JOIN users ON posts.id_users = users.id
    INNER JOIN post_category ON posts.id = post_category.post_id
    INNER JOIN category ON post_category.catego_id = category.id
    LEFT JOIN (
		select post_id, count(*) as comment_count
		from comments
		GROUP by post_id
		) as comment on comment.post_id = posts.id
    GROUP BY posts.id
    ORDER BY posts.created_at DESC;
	`)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return nil, err
	}
	return rows, nil
}

func Loggedin(w http.ResponseWriter, r *http.Request) bool {
	if db.DB == nil {
		log.Println("Database connection is nil!")
		return false
	}
	cookie, err := r.Cookie("userId")
	if err != nil || cookie == nil {
		return false
	}

	// Just check if session exists, don't trigger logout
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE uuid = ?)", cookie.Value).Scan(&exists)
	if err != nil {
		log.Println("Error checking session:", err)
		return false
	}
	return exists
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "templates/inde.html") TODO: fix
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}

func Islogged(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := Loggedin(w, r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"isLoggedIn": isLoggedIn,
	})
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := Loggedin(w, r)
	if r.Method == http.MethodPost {
		return
	}
	type CategoryData struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		PostCount int    `json:"postCount"`
	}
	type PageData struct {
		Categories []CategoryData `json:"categories"`
		IsLoggedIn bool           `json:"isLoggedIn"`
	}

	categories := []CategoryData{}

	rows, err := db.DB.Query(`
        SELECT c.id, c.name, COUNT(pc.post_id) AS post_count 
        FROM category c
        LEFT JOIN post_category pc ON c.id = pc.catego_id
        GROUP BY c.id, c.name
    `)
	if err != nil {
		http.Error(w, `{"error": "Unable to fetch categories"}`, http.StatusInternalServerError)
		log.Println("Error fetching categories:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category CategoryData
		if err := rows.Scan(&category.ID, &category.Name, &category.PostCount); err != nil {
			http.Error(w, `{"error": "Unable to process category data"}`, http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}
	data := PageData{
		Categories: categories,
		IsLoggedIn: isLoggedIn,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func errorsHandler(w http.ResponseWriter, r *http.Request) {
	islogedin := Loggedin(w, r)
	code := r.URL.Query().Get("code")
	errorin := map[string]string{
		"404": "Path Not Found",
		"401": "Not Allowed",
	}

	if code == "404" {
		w.WriteHeader(http.StatusNotFound)
	} else if code == "401" {
		w.WriteHeader(http.StatusUnauthorized)
	}
	data := struct {
		IsLoggedIn bool
		Title      string
		Message    string
	}{
		IsLoggedIn: islogedin,
		Title:      code,
		Message:    errorin[code],
	}
	tmpl := template.Must(template.ParseFiles("templates/errors.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := Loggedin(w, r)
	if !isLoggedIn {
		http.Redirect(w, r, "/errors?code=401", http.StatusFound) //////////////////////////////////////////////////////////////////////////////////////////////////////
		return
	}
	// Wrap the isLoggedIn value in a struct
	data := struct {
		IsLoggedIn bool
	}{
		IsLoggedIn: isLoggedIn,
	}
	// Pass the struct to the template
	tmpl := template.Must(template.ParseFiles("templates/new_post.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func Fetchuid(token string) (int, error) {
	var id string
	var iid int
	err := db.DB.QueryRow("SELECT id FROM users WHERE uuid = ?", token).Scan(&id)
	if err != nil {
		return 0, err
	}
	iid, err = strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return iid, nil
}

func Getusername(id int) (string, error) {
	var usrname string
	err := db.DB.QueryRow("select username from users where id = ?", id).Scan(&usrname)
	if err != nil {
		return "", err
	}
	return usrname, nil
}

func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid method"}`, http.StatusMethodNotAllowed)
		return
	}

	isLoggedIn := Loggedin(w, r)
	cookie, err := r.Cookie("userId")
	if err != nil || !isLoggedIn {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	userID, err := Fetchuid(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	var post struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Catego     []string `json:"category"`
		StatusPost string   `json:"statusPost"`
		Attachment string   `json:"attachment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	fmt.Println(post.Attachment)
	whitespaceRegex := regexp.MustCompile(`^\s*$`)
	if post.Title == "" || post.Content == "" || whitespaceRegex.MatchString(post.Title) || whitespaceRegex.MatchString(post.Content) {
		http.Error(w, `{"error": "Title and content are required"}`, http.StatusBadRequest)
		return
	}
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	// Insert post
	result, err := tx.Exec("INSERT INTO posts (id_users, title, content, attachement , status) VALUES (?, ?, ?,?,?)",
		userID, post.Title, post.Content, post.Attachment, post.StatusPost)
	if err != nil {
		http.Error(w, `{"error": "Error saving post"}`, http.StatusInternalServerError)
		return
	}
	postID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, `{"error": "Error retrieving post ID"}`, http.StatusInternalServerError)
		return
	}
	for _, cname := range post.Catego {
		var categoryID int64

		err = tx.QueryRow("SELECT id FROM category WHERE name = ?", cname).Scan(&categoryID)
		if err != nil {
			if err == sql.ErrNoRows {

				results, err := tx.Exec("INSERT INTO category (name) VALUES (?)", cname)
				if err != nil {
					http.Error(w, "Error creating category", http.StatusInternalServerError)
					return
				}

				categoryID, err = results.LastInsertId()
				if err != nil {
					http.Error(w, "Error retrieving category ID", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Error fetching category", http.StatusInternalServerError)
				return
			}
		}

		_, err = tx.Exec("INSERT INTO post_category (catego_id, post_id) VALUES (?, ?)", categoryID, postID)
		if err != nil {
			http.Error(w, "Error saving post category relation", http.StatusInternalServerError)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	hub.HUB.Brdcast <- []byte(`{"type": "new_post", "islogged": ` + strconv.FormatBool(isLoggedIn) + `}`)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
	})
}

// func isValidEmail(email string) bool {
// 	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
// 	re := regexp.MustCompile(emailRegex)
// 	return re.MatchString(email)
// }

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if Loggedin(w, r) {
		http.Error(w, `{"error": "Already logged in"}`, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	// log.Println("Request body:", string(body))

	var user struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       string `json:"age"`
		Gender    string `json:"gender"`
		Fname     string `json:"fname"`
		Lname     string `json:"lname"`
		Birthdate string `json:"birthdate"`
		Avatar    string `json:"avatar"`
		Aboutme   string `json:"aboutme"`
		Status    string `json:"status"`
	}

	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	// fmt.Println(user)
	// Validate input
	containSpaceRegex := regexp.MustCompile(`\s`)
	whitespaceRegex := regexp.MustCompile(`^\s*$`)
	if user.Username == "" || user.Email == "" || user.Password == "" ||
		containSpaceRegex.MatchString(user.Username) || containSpaceRegex.MatchString(user.Email) ||
		whitespaceRegex.MatchString(user.Password) {
		log.Println("Validation failed for user input")
		http.Error(w, `{"error": "Username, email, or password should not be empty or contain spaces"}`, http.StatusBadRequest)
		return
	}

	// Check if user exists
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?)", user.Username, user.Email).Scan(&exists)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	if exists {
		log.Println("User already exists:", user.Username, user.Email)
		http.Error(w, `{"error": "Username or email already exists"}`, http.StatusBadRequest)
		return
	}

	// Insert user
	uuid, err := generateToken()
	if err != nil {
		log.Println("Error generating token:", err)
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec(`
    INSERT INTO users (
        uuid, username, email, birthdate, password, age, gender, fname, lname, avatar, aboutme, status
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)`,
		uuid, user.Username, user.Email, user.Birthdate, user.Password, user.Age, user.Gender, user.Fname, user.Lname, user.Avatar, user.Aboutme, user.Status)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "User registered successfully"}

	json.NewEncoder(w).Encode(response)
}

// func registerHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	if Loggedin(w, r) {
// 		http.Error(w, `{"error": "Already logged in"}`, http.StatusBadRequest)
// 		return
// 	}

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		log.Println("Error reading request body:", err)
// 		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
// 		return
// 	}
// 	// log.Println("Request body:", string(body))

// 	var user struct {
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 		Age      string `json:"age"`
// 		Gender   string `json:"gender"`
// 		Fname    string `json:"fname"`
// 		Lname    string `json:"lname"`
// 	}
// 	if err := json.Unmarshal(body, &user); err != nil {
// 		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
// 		return
// 	}
// 	// fmt.Println(user)
// 	// Validate input
// 	containSpaceRegex := regexp.MustCompile(`\s`)
// 	whitespaceRegex := regexp.MustCompile(`^\s*$`)
// 	if user.Username == "" || user.Email == "" || user.Password == "" ||
// 		containSpaceRegex.MatchString(user.Username) || containSpaceRegex.MatchString(user.Email) ||
// 		whitespaceRegex.MatchString(user.Password) {
// 		log.Println("Validation failed for user input")
// 		http.Error(w, `{"error": "Username, email, or password should not be empty or contain spaces"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Check if user exists
// 	var exists bool
// 	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?)", user.Username, user.Email).Scan(&exists)
// 	if err != nil {
// 		log.Println("Error checking if user exists:", err)
// 		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
// 		return
// 	}
// 	if exists {
// 		log.Println("User already exists:", user.Username, user.Email)
// 		http.Error(w, `{"error": "Username or email already exists"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Insert user
// 	uuid, err := generateToken()
// 	if err != nil {
// 		log.Println("Error generating token:", err)
// 		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	_, err = db.DB.Exec("INSERT INTO users (uuid, username, email, password, age, gender, fname, lname) VALUES (?, ?, ?, ?,?,?,?,?)", uuid, user.Username, user.Email, user.Password, user.Age, user.Gender, user.Fname, user.Lname)
// 	if err != nil {
// 		log.Println("Error inserting user into database:", err)
// 		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]string{"message": "User registered successfully"}

// 	json.NewEncoder(w).Encode(response)
// }

// func setandchecksession(cokie *http.Cookie) bool {
// 	var exists bool
// 	err := db.DB.QueryRow(
// 		"SELECT EXISTS(SELECT 1 FROM users WHERE uuid = ?)",
// 		cokie.Value,
// 	).Scan(&exists)
// 	if err != nil {
// 		log.Println("Error checking existing user:", err)
// 		return false
// 	}
// 	if exists {
// 		return false
// 	}
// 	return true
// }

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	bytes[6] = (bytes[6] & 0x0f) | 0x40 // UUID version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // UUID variant 1

	return hex.EncodeToString(bytes), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	var storedPassword string
	var id string
	var usrnm string
	err := db.DB.QueryRow("SELECT password, id FROM users WHERE username = ?", credentials.Username).Scan(&storedPassword, &id)
	if err != nil || storedPassword != credentials.Password {
		err = db.DB.QueryRow("SELECT password, id FROM users WHERE email = ?", credentials.Username).Scan(&storedPassword, &id)
		if err != nil || storedPassword != credentials.Password {
			http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
			return
		} else {
			err = db.DB.QueryRow("SELECT username  FROM users WHERE email = ?", credentials.Username).Scan(&usrnm)
			if err != nil {
				http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
				return
			}
			credentials.Username = usrnm
		}
	}

	err = db.DB.QueryRow("SELECT password, id FROM users WHERE username = ? OR email = ?", credentials.Username, credentials.Username).Scan(&storedPassword, &id)
	if err != nil || storedPassword != credentials.Password {
		http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	// Generate and update UUID
	uuid, err := generateToken()
	if err != nil {
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	_, err = db.DB.Exec("UPDATE users SET uuid = ? WHERE id = ?", uuid, id)
	if err != nil {
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	// Set cookie
	cookie := &http.Cookie{
		Name:     "userId",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(168 * time.Hour),
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login successful",
		"username": credentials.Username,
	})
}

// lougoutiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userId")
	if err != nil {
		return
	}

	// fmt.Println("yup")
	var username string
	err = db.DB.QueryRow("SELECT username FROM users WHERE uuid = ?", cookie.Value).Scan(&username)
	// fmt.Println(username)
	if err != nil {
		http.Error(w, `{"error": "Invalid username"}`, http.StatusUnauthorized)
		return
	}
	// if client, ok := hub.srsu[username]; ok {
	//     hub.unregister <- client
	// }
	cookie = &http.Cookie{
		Name:     "userId",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	hub.HUB.Unregister <- hub.HUB.Srsu[username]
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/css/" || strings.HasSuffix(r.URL.Path, "/") {
		http.Error(w, "YOU are not allowed", http.StatusUnauthorized)
		return
	}
	filePath := "./css" + strings.TrimPrefix(r.URL.Path, "/css")
	http.ServeFile(w, r, filePath)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	profileuser := r.URL.Query().Get("user")
	isLoggedIn := Loggedin(w, r)
	if !isLoggedIn {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "You are not loggedd hiiiiin",
		})
		return
	}
	cookie, err := r.Cookie("userId")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid cookie",
		})
		return
	}

	userId, err := Fetchuid(cookie.Value)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid user ID",
		})
		return
	}

	var user structs.User
	var posts []structs.Post

	err = db.DB.QueryRow("SELECT id, username, email , fname, lname,status FROM users WHERE username = ?", profileuser).
		Scan(&user.ID, &user.Username, &user.Email, &user.Fname, &user.Lname, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "User not found",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Internal Server Error",
		})
		return
	}
	rows, err := db.DB.Query("SELECT id, title FROM posts WHERE id_users = ?", user.ID)
	if err != nil {
	}
	for rows.Next() {
		var pos structs.Post
		if err := rows.Scan(&pos.ID, &pos.Title); err != nil {
			log.Println("Error scanning post:", err)
			continue
		}
		posts = append(posts, pos)
	}
	rows, err = db.DB.Query("SELECT follower , followed FROM followers WHERE follower = ? or followed = ? ", user.ID, user.ID)
	if err != nil {
	}
	for rows.Next() {
		fd := 0
		fr := 0
		if err := rows.Scan(&fr, &fd); err != nil {
		}
		if fd == user.ID {
			user.Followers = append(user.Followers, fr)
		} else if fr == user.ID {
			user.Followed = append(user.Followed, fd)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	if user.ID == userId {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"username":   user.Username,
			"email":      user.Email,
			"isLoggedIn": isLoggedIn,
			"fname":      user.Fname,
			"lname":      user.Lname,
			"posts":      posts,
			"followers":  user.Followers,
			"followed":   user.Followed,
		})
	} else {
		checkfd := false
		for rng := range user.Followers {
			if rng == userId {
				checkfd = true
				break
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"username":   user.Username,
			"status":     user.Status,
			"isLoggedIn": isLoggedIn,
			"fname":      user.Fname,
			"lname":      user.Lname,
			"posts":      posts,
			"followers":  user.Followers,
			"followed":   user.Followed,
			"unfollow":   checkfd,
		})
	}
}

func Forunf(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hiiiiiiiiiiiiiiii")
	profileuser := r.URL.Query().Get("id")
	isLoggedIn := Loggedin(w, r)
	if !isLoggedIn {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "You are not logged in",
		})
		return
	}
	cookie, err := r.Cookie("userId")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid cookie",
		})
		return
	}
	fmt.Println("looooooooooooooooooooooooooooooool", profileuser)

	userId, err := Fetchuid(cookie.Value)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid user ID",
		})
		return
	}
	var me string
	var find string
	err = db.DB.QueryRow("select username from users where id = ?", userId).Scan(&me)
	if err != nil {
	}
	err = db.DB.QueryRow("SELECT follower FROM followers WHERE follower = ? AND followed = ? ", me, profileuser).Scan(&find)
	if err != nil {
	}
	if find == me {
		_, err := db.DB.Exec("DELETE FROM followers WHERE follower = ? AND followed = ?", me, profileuser)
		if err != nil {
		}
		find = "unfollowed"
	} else {
		var stat string
		err = db.DB.QueryRow("select status from users where username = ?", profileuser).Scan(&stat)
		if err != nil {
		}
		if stat == "public" {
			_, err := db.DB.Exec("INSERT INTO followers (follower ,followed) VALUES(?,?)", me, profileuser)
			fmt.Println("looooooooooooooooooooooooooooooool")
			if err != nil {
			}
			find = "followed"
			_, err = db.DB.Exec("insert into notifications (notif,user) values(?,?)", me+" start to follow you", profileuser)
			if err != nil {
			}
		} else {
			_, err = db.DB.Exec("insert into notifications (notif,user) values(?,?)", me+" send you follow request", profileuser)
			if err != nil {
			}
			find = "follow request sent"
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": find,
	})
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := Loggedin(w, r)
	if !isLoggedIn {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie("userId")
	if err != nil {
		http.Error(w, `{"error": "Invalid cookie"}`, http.StatusBadRequest)
		return
	}

	userID, err := Fetchuid(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	otherUser := r.URL.Query().Get("user")
	if otherUser == "" {
		http.Error(w, `{"error": "User parameter is required"}`, http.StatusBadRequest)
		return
	}

	var otherUserID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = ?", otherUser).Scan(&otherUserID)
	if err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	rows, err := db.DB.Query(`
        SELECT m.content, m.created_at, u.username 
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE (m.sender_id = ? AND m.receiver_id = ?)
           OR (m.sender_id = ? AND m.receiver_id = ?)
        ORDER BY m.created_at ASC`,
		userID, otherUserID, otherUserID, userID)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var content, createdAt, username string
		if err := rows.Scan(&content, &createdAt, &username); err != nil {
			continue
		}
		messages = append(messages, map[string]interface{}{
			"sender":  username,
			"content": content,
			"time":    createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
	})
}

func CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		http.Error(w, "Category ID is missing", http.StatusBadRequest)
		return
	}
	// fmt.Println(categoryID,"looool")
	// Check if the user is logged in
	IsLoggedIn := Loggedin(w, r)

	// type PostData struct {
	// 	ID           int
	// 	Title        string
	// 	CreatedAt    string
	// 	Username     string
	// 	LikeCount    int
	// 	DislikeCount int
	// 	CommentCount int
	// }

	// posts := []PostData{}

	query := `
        SELECT 
            p.id,
            p.title, 
            p.created_at, 
            u.username,
			COALESCE(comment.comment_count, 0) AS comment
        FROM posts p
        JOIN users u ON p.id_users = u.id
        JOIN post_category pc ON p.id = pc.post_id
		 LEFT JOIN (
		select post_id, count(*) as comment_count
		from comments
		GROUP by post_id
		) as comment on comment.post_id = p.id
        WHERE pc.catego_id = ?
        GROUP BY p.id, u.username
    `
	rows, err := db.DB.Query(query, categoryID)
	if err != nil {
		http.Error(w, "Unable to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.CreatedAt, &post.Username, &post.CommentCount); err != nil {
			log.Println("Error scanning post:", err)
			continue
		}
		posts = append(posts, post)
	}

	// tmpl, err := template.ParseFiles("templates/category_posts.html")
	// if err != nil {
	// 	http.Error(w, "Unable to load category posts template", http.StatusInternalServerError)
	// 	return
	// }
	var catego string
	db.DB.QueryRow("select name from category where id = ?", &categoryID).Scan(&catego)

	// tmpl.Execute(w, map[string]interface{}{
	// 	"Posts":      posts,
	// 	"IsLoggedIn": IsLoggedIn,
	// 	"Categorie": catego,
	// })
	// fmt.Println(posts,IsLoggedIn,catego)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"posts":      posts,
		"isLoggedIn": IsLoggedIn,
		"catego":     catego,
	})
}

func ConversationsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userId")
	if err != nil {
		http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
		return
	}

	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE uuid = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, `{"error": "Invalid user"}`, http.StatusBadRequest)
		return
	}
	query := `
    WITH user_conversations AS (
        SELECT DISTINCT 
            CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as partner_id
        FROM messages
        WHERE sender_id = ? OR receiver_id = ?
    )
    SELECT 
        u.username,
        MAX(m.created_at) as last_message_time,
        COUNT(CASE WHEN m.receiver_id = ? AND m.read = 0 THEN 1 END) as unread_count,
        last_msg.content as last_message_content
    FROM user_conversations uc
    JOIN users u ON uc.partner_id = u.id
    LEFT JOIN messages m ON (m.sender_id = uc.partner_id AND m.receiver_id = ?) 
                        OR (m.sender_id = ? AND m.receiver_id = uc.partner_id)
    LEFT JOIN messages last_msg ON last_msg.id = (
        SELECT id FROM messages 
        WHERE (sender_id = uc.partner_id AND receiver_id = ?)
           OR (sender_id = ? AND receiver_id = uc.partner_id)
        ORDER BY created_at DESC LIMIT 1
    )
    GROUP BY u.username, last_msg.content
    ORDER BY last_message_time DESC`

	rows, err := db.DB.Query(query,
		userID, userID, userID, userID, userID, userID, userID, userID)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Conversation struct {
		User        string  `json:"user"`
		LastMessage *string `json:"last_message"`
		LastMsgTim  *string `json:"msg_time"`
		Unread      int     `json:"unread"`
		Enligne     bool    `json:"enligne"`
	}
	type conversationandusers struct {
		Conv  []Conversation `json:"conv"`
		Users []string       `json:"users"`
	}
	var retrn conversationandusers

	retrn.Users, err = getUserNames()
	// fmt.Println(retrn.Users);
	// var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var lastMsgTime sql.NullString
		var lastMsgContent sql.NullString

		if err := rows.Scan(&conv.User, &lastMsgTime, &conv.Unread, &lastMsgContent); err != nil {
			log.Printf("Error scanning conversation row: %v", err)
			continue
		}
		if lastMsgContent.Valid {
			conv.LastMessage = &lastMsgContent.String
			conv.LastMsgTim = &lastMsgTime.String
		}
		conv.Enligne = false
		retrn.Conv = append(retrn.Conv, conv)
	}

	if len(retrn.Conv) == 0 {
		retrn.Conv = []Conversation{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(retrn); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func getUserNames() ([]string, error) {
	rows, err := db.DB.Query("SELECT username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		userNames = append(userNames, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userNames, nil
}

// func conversationsHandler(w http.ResponseWriter, r *http.Request) {
//     cookie, err := r.Cookie("userId")
//     if err != nil {
//         http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
//         return
//     }

//     var userID int
//     err = db.DB.QueryRow("SELECT id FROM users WHERE uuid = ?", cookie.Value).Scan(&userID)
//     if err != nil {
//         http.Error(w, `{"error": "Invalid user"}`, http.StatusBadRequest)
//         return
//     }

//     rows, err := db.DB.Query(`
//         SELECT
//             u.username,
//             MAX(m.created_at) as last_message_time,
//             COUNT(CASE WHEN m.receiver_id = ? AND m.read = 0 THEN 1 END) as unread_count,
//             last_msg.content as last_message_content
//         FROM (
//             SELECT DISTINCT
//                 CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as partner_id
//             FROM messages
//             WHERE sender_id = ? OR receiver_id = ?
//         ) as partners
//         JOIN users u ON partners.partner_id = u.id
//         LEFT JOIN messages m ON (m.sender_id = partners.partner_id AND m.receiver_id = ?)
//                             OR (m.sender_id = ? AND m.receiver_id = partners.partner_id)
//         LEFT JOIN messages last_msg ON last_msg.id = (
//             SELECT id FROM messages
//             WHERE (sender_id = partners.partner_id AND receiver_id = ?)
//                OR (sender_id = ? AND receiver_id = partners.partner_id)
//             ORDER BY created_at DESC LIMIT 1
//         )
//         GROUP BY u.username, last_msg.content
//         ORDER BY last_message_time DESC`,
//         userID, userID, userID, userID, userID, userID, userID, userID)
//     if err != nil {
//         http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
//         return
//     }
//     defer rows.Close()

//     type Conversation struct {
//         User        string     `json:"user"`
//         LastMessage *string    `json:"last_message"`
//         Unread      int        `json:"unread"`
//     }

//     var conversations = make([]Conversation,0)
//     for rows.Next() {
//         var conv Conversation
//         var lastMsgTime sql.NullString
//         var lastMsgContent sql.NullString
//         if err := rows.Scan(&conv.User, &lastMsgTime, &conv.Unread, &lastMsgContent); err != nil {
//             continue
//         }
//         if lastMsgContent.Valid {
// 			fmt.Println("lol")
//             conv.LastMessage = &lastMsgContent.String
//         }
//         conversations = append(conversations, conv)
//     }
// 	fmt.Println(conversations[0].LastMessage.)
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(conversations)
// }

// func conversationsHandler(w http.ResponseWriter, r *http.Request) {
//     // Authentication check
//     cookie, err := r.Cookie("userId")
//     if err != nil {
//         http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
//         return
//     }

//     var userID int
//     err = db.DB.QueryRow("SELECT id FROM users WHERE uuid = ?", cookie.Value).Scan(&userID)
//     if err != nil {
//         http.Error(w, `{"error": "Invalid user"}`, http.StatusBadRequest)
//         return
//     }

//     // Get distinct conversation partners
//     rows, err := db.DB.Query(`
//         SELECT DISTINCT
//             u.id,
//             u.username,
//             (SELECT COUNT(*) FROM messages
//              WHERE sender_id = u.id AND receiver_id = ? AND read = 0) AS unread_count
//         FROM messages m
//         JOIN users u ON (m.sender_id = u.id OR m.receiver_id = u.id) AND u.id != ?
//         WHERE m.sender_id = ? OR m.receiver_id = ?
//         ORDER BY u.username`,
//         userID, userID, userID, userID)
//     if err != nil {
// 		fmt.Println(err)
//         http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
//         return
//     }
//     defer rows.Close()

//     type Conversation struct {
//         UserID   int     `json:"user_id"`
//         Username string  `json:"username"`
//         Unread   int     `json:"unread"`
//     }

//     var conversations = make([]Conversation,0)
//     for rows.Next() {
//         var conv Conversation
//         if err := rows.Scan(&conv.UserID, &conv.Username, &conv.Unread); err != nil {
//             continue
//         }
//         conversations = append(conversations, conv)
//     }
// 	fmt.Println(conversations,"yes")
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(conversations)
// }

func MarkReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("userId")
	if err != nil {
		http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
		return
	}

	otherUser := r.URL.Query().Get("user")
	if otherUser == "" {
		http.Error(w, `{"error": "User parameter is required"}`, http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
        UPDATE messages SET read = 1 
        WHERE sender_id = (SELECT id FROM users WHERE username = ?)
        AND receiver_id = (SELECT id FROM users WHERE uuid = ?)`,
		otherUser, cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
