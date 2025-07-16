package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/logs"
	"social-network/app/modules"
	"social-network/app/structs"
)

type Profile struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DisplayName    string `json:"display_name"`
	DateOfBirth    string `json:"date_of_birth,omitempty"`
	Gender         string `json:"gender"`
	Avatar         string `json:"avatar"`
	Description    string `json:"description"`
	IsPublic       bool   `json:"isPublic"`
	IsUser         bool   `json:"isUser"`
	CreatedAt      string `json:"created_at"`
	IsSelf         bool   `json:"isSelf"`
	FollowStatus   string `json:"followStatus"`
	PostCount      int    `json:"post_count"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	nickname := strings.TrimSpace(r.PathValue("name"))
	payload, ok := r.Context().Value(auth.UserContextKey).(*jwt.JwtPayload)

	var profile Profile
	var temp sql.NullString
	var avatarTmp sql.NullString
	var err error
	if nickname == "me" {
		nickname = payload.Username
	}
	// 👤 Case 1: viewer requests their own profile using their nickname
	if ok && strings.EqualFold(nickname, payload.Username) {
		// Fetch by ID (self profile)
		err = modules.DB.QueryRow(`
			SELECT
			    pr.id,
			    pr.email,
			    pr.first_name,
			    pr.last_name,
			    pr.display_name,
			    pr.date_of_birth,
			    pr.gender,
			    pr.avatar,
			    pr.description,
			    pr.is_public,
			    pr.is_user,
			    pr.created_at,
			    (
			        SELECT
			            COUNT(*)
			        FROM
			            posts p
			        WHERE
			            p.user_id = pr.id
			            AND p.group_id IS NULL
			    ) AS post_count,
			    (
			        SELECT
			            COUNT(*)
			        FROM
			            follow f
			        WHERE
			            f.following_id = pr.id
			    ) AS follower_count,
			    (
			        SELECT
			            COUNT(*)
			        FROM
			            follow f
			        WHERE
			            f.follower_id = pr.id
			    ) AS following_count
			FROM
			    profile pr
			WHERE
			    id = ?;
		`, payload.Sub).Scan(
			&profile.ID,
			&profile.Email,
			&profile.FirstName,
			&profile.LastName,
			&profile.DisplayName,
			&profile.DateOfBirth,
			&profile.Gender,
			&avatarTmp,
			&temp,
			&profile.IsPublic,
			&profile.IsUser,
			&profile.CreatedAt,
			&profile.PostCount,
			&profile.FollowerCount,
			&profile.FollowingCount,
		)

		if avatarTmp.Valid {
			profile.Avatar = avatarTmp.String
		}
		if temp.Valid {
			profile.Description = temp.String
		}
		profile.IsSelf = true

	} else {
		// 🕵️‍♂️ Case 2: someone else’s profile
		err = modules.DB.QueryRow(`
			SELECT ID,email, first_name, last_name, display_name, date_of_birth, gender,
			       avatar, description, is_public, is_user, created_at,
				   (SELECT COUNT(*) FROM posts p WHERE p.user_id = ID) AS post_count,
				   (SELECT COUNT(*) FROM follow f WHERE f.following_id = ID) AS follower_count,
				   (SELECT COUNT(*) FROM follow f WHERE f.follower_id = ID) AS following_count
			FROM profile
			WHERE LOWER(display_name) = LOWER(?)
		`, nickname).Scan(
			&profile.ID,
			&profile.Email,
			&profile.FirstName,
			&profile.LastName,
			&profile.DisplayName,
			&profile.DateOfBirth,
			&profile.Gender,
			&avatarTmp,
			&temp,
			&profile.IsPublic,
			&profile.IsUser,
			&profile.CreatedAt,
			&profile.PostCount,
			&profile.FollowerCount,
			&profile.FollowingCount,
		)
		if err != nil {
			logs.ErrorLog.Println("DB error:", err)
			auth.JsRespond(w, "Bad request", http.StatusBadRequest)
			return
		}

		if profile.IsUser && !profile.IsPublic {
			profile.Email = ""
			profile.FirstName = ""
			profile.LastName = ""
			profile.DateOfBirth = ""
			profile.Gender = ""
		}

		if temp.Valid {
			profile.Description = temp.String
		}

		if avatarTmp.Valid {
			profile.Avatar = avatarTmp.String
		}

		profile.FollowStatus, err = modules.GetRelationship(payload.Sub, profile.ID)
		if err != nil {
			auth.JsRespond(w, "Feild to get relationship", http.StatusNotFound)
			return
		}

		profile.IsSelf = false
	}

	if err != nil {
		if err == sql.ErrNoRows {
			auth.JsRespond(w, "Profile not found", http.StatusNotFound)
		} else {
			logs.ErrorLog.Println("DB error:", err)
			auth.JsRespond(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logs.ErrorLog.Println("Invalid userId:", err)
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	rows, err := modules.DB.Query(`
		SELECT
		    f.follower_id,
		    p.display_name,
		    p.avatar
		FROM
		    follow f
		    JOIN profile p ON (
		        f.follower_id = p.id
		        AND p.is_user = 1
		    )
		WHERE
		    f.following_id = ?
	`, userId)
	if err != nil {
		logs.ErrorLog.Printf("GetFollowers query error: %q", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followers []structs.UsersGet
	for rows.Next() {
		var follower structs.UsersGet
		if err := rows.Scan(&follower.ID, &follower.Username, &follower.Avatar); err != nil {
			logs.ErrorLog.Printf("Error scanning followers: %q", err.Error())
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		followers = append(followers, follower)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logs.ErrorLog.Println("Invalid userId:", err)
		auth.JsRespond(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	rows, err := modules.DB.Query(`
		SELECT
		    f.following_id,
		    p.display_name,
		    p.avatar
		FROM
		    follow f
		    JOIN profile p ON f.following_id = p.id
		    AND p.is_user = 1
		WHERE
		    f.follower_id = ?;`, userId)
	if err != nil {
		logs.ErrorLog.Printf("GetFollowing query error: %q", err.Error())
		json.NewEncoder(w).Encode(err)
		return
	}
	defer rows.Close()

	var following []structs.UsersGet
	for rows.Next() {
		var follow structs.UsersGet
		if err := rows.Scan(&follow.ID, &follow.Username, &follow.Avatar); err != nil {
			logs.ErrorLog.Printf("Error scanning following: %q", err.Error())
			json.NewEncoder(w).Encode(err)
			return
		}
		following = append(following, follow)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}
