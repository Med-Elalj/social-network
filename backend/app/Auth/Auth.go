package auth

import (
	"encoding/json"
	"net/http"

	"social-network/app/Auth/jwt"
	"social-network/app/logs"
	"social-network/app/modules"
)

func Authorize(w http.ResponseWriter, r *http.Request, userID int) {
	username, err := GetElemVal[string]("display_name", "profile", `id = ?`, userID)
	if err != nil {
		logs.ErrorLog.Println("Error getting username:", err)
		JsRespond(w, "Bad request", http.StatusBadRequest)
		return
	}

	jwtToken, sessionID, refreshToken, err := CheckSession(r, userID, username)
	if err != nil {
		logs.ErrorLog.Println(err)
		JsRespond(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Access Token
	SetCookie(w, AuthInfo.JwtTokenName, jwtToken, int(AuthInfo.JwtExpiration.Seconds()))
	// Session ID & Refresh Token (HttpOnly)
	SetCookie(w, AuthInfo.SessionIDName, sessionID, int(AuthInfo.SessionExpiration.Seconds()))
	// Refresh Token
	SetCookie(w, AuthInfo.RefreshTokenName, refreshToken, int(AuthInfo.SessionExpiration.Seconds()))
}

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Get cookies
	jwtCookie, errJWT := r.Cookie(AuthInfo.JwtTokenName)
	sidCookie, errSID := r.Cookie(AuthInfo.SessionIDName)

	if errJWT != nil || errSID != nil {
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	jwtToken := jwtCookie.Value
	sessionID := sidCookie.Value

	// 2. Validate JWT and Session
	payload, err := jwt.JWTVerify(jwtToken)
	if err != nil || payload == nil {
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	// 3. Check if session is still active in DB
	validSession, err := SessionExists(payload.Sub, sessionID)
	if err != nil || !validSession || sessionID != payload.SessionID {
		ClearCookie(w, AuthInfo.SessionIDName)
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}
	// 4. Check if IP and User-Agent match
	Session, _ := GetSessionByID(sidCookie.Value)
	clientIP := GetIP(r)
	if Session.IP != clientIP || Session.UserAgent != r.UserAgent() {
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	// 5. Respond with success
	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}

// Logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// delete from db
	sidCookie, _ := r.Cookie(AuthInfo.SessionIDName)
	rtCookie, _ := r.Cookie(AuthInfo.RefreshTokenName)
	// Invalidate session in DB
	_, err := modules.DB.Exec(`DELETE FROM sessions WHERE session_id = ? AND refresh_token = ?`, sidCookie.Value, rtCookie.Value)
	if err != nil {
		logs.ErrorLog.Println("Error deleting session:", err)
		JsRespond(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Clear session cookies
	ClearCookie(w, AuthInfo.JwtTokenName)
	ClearCookie(w, AuthInfo.SessionIDName)
	ClearCookie(w, AuthInfo.RefreshTokenName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
