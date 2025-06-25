package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social-network/app/Auth/jwt"
)

type contextKey string

const UserContextKey contextKey = "user"

func authorize(w http.ResponseWriter, userID int) {
	// username := getElemVal("username", "users", `id = "`+strconv.Itoa(userID)+`"`).(string)
	// jwt, sessionID, err := CheckSession(userID, username)
	// if err != nil {
	// 	helpers.ErrorLog.Println(err)
	// }

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "jwt",
	// 	Value:    jwt,
	// 	Path:     "/",
	// 	Secure:   true,
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteStrictMode,
	// 	Expires:  time.Now().Add(60 * time.Minute),
	// })

	// // Set the Session ID in a separate HttpOnly cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "ssid",
	// 	Value:    sessionID,
	// 	Path:     "/",
	// 	Secure:   true,
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteStrictMode,
	// 	Expires:  time.Now().Add(60 * time.Minute),
	// })

	// w.WriteHeader(http.StatusOK)
}

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	// // check online from js based on the api
	// jwt_token, _ := ExtractJWT(r)
	// ssid, _ := ExtractSSID(r)
	// auth, err := VerifyUser(jwt_token, ssid)
	// count, _ := helpers.EntryExists("session_id", ssid, "sessions", false)
	// if count != 1 {
	// 	http.SetCookie(w, &http.Cookie{
	// 		Name:    "ssid",
	// 		Value:   "",
	// 		Path:    "/",
	// 		MaxAge:  -1,
	// 		Expires: time.Unix(0, 0),
	// 	})
	// }

	// if !auth || err != nil {
	// 	json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
	// 	return
	// }

	// json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}

// Check if user logged
func Islogged(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := r.Context().Value(UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "Unauthorized"}`)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": data.Username + " User is logged in",
	})
}

// Logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:     "JWT",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,                   // ensures deletion
		HttpOnly: true,                 // match original
		Secure:   true,                 // match original
		SameSite: http.SameSiteLaxMode, // match original
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
