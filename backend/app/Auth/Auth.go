package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"social-network/app/Auth/jwt"
	"social-network/app/modules"
	"social-network/server/logs"
)

func authorize(w http.ResponseWriter, r *http.Request, userID int) {
	username := GetElemVal("display_name", "profile", `id = "`+strconv.Itoa(userID)+`"`).(string)

	jwtToken, sessionID, refreshToken, err := CheckSession(r, userID, username)
	if err != nil {
		logs.ErrorLog.Println(err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Access Token
	SetCookie(w, AuthInfo.JwtTokenName, jwtToken, int(AuthInfo.JwtExpiration.Seconds()))
	// Session ID & Refresh Token (HttpOnly)
	SetCookie(w, AuthInfo.SessionIDName, sessionID, int(AuthInfo.SessionExpiration.Seconds()))
	// Refresh Token
	SetCookie(w, AuthInfo.RefreshTokenName, refreshToken, int(AuthInfo.SessionExpiration.Seconds()))

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
}

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Get cookies
	jwtCookie, errJWT := r.Cookie(AuthInfo.JwtTokenName)
	sidCookie, errSID := r.Cookie(AuthInfo.SessionIDName)

	if errJWT != nil || errSID != nil {
		ClearCookie(w, AuthInfo.JwtTokenName)
		ClearCookie(w, AuthInfo.SessionIDName)
		ClearCookie(w, AuthInfo.RefreshTokenName)
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	jwtToken := jwtCookie.Value
	sessionID := sidCookie.Value

	// 2. Validate JWT and Session
	payload, err := jwt.JWTVerify(jwtToken)
	if err != nil || payload == nil {
		ClearCookie(w, AuthInfo.JwtTokenName)
		ClearCookie(w, AuthInfo.SessionIDName)
		ClearCookie(w, AuthInfo.RefreshTokenName)
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	// 3. Check if session is still active in DB
	validSession, err := SessionExists(payload.Sub, sessionID)
	if err != nil || !validSession {
		ClearCookie(w, AuthInfo.JwtTokenName)
		ClearCookie(w, AuthInfo.SessionIDName)
		ClearCookie(w, AuthInfo.RefreshTokenName)
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	// 4. Optionally: set user in context
	ctx := context.WithValue(r.Context(), UserContextKey, payload)
	r = r.WithContext(ctx)

	// 5. Respond with success
	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}

// Logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// delete from db
	sidCookie, errSID := r.Cookie(AuthInfo.SessionIDName)
	if errSID != nil {
		http.Error(w, "Unauthorized - missing cookies", http.StatusUnauthorized)
		return
	}
	rtCookie, errRT := r.Cookie(AuthInfo.RefreshTokenName)
	if errRT != nil {
		http.Error(w, "Unauthorized - missing cookies", http.StatusUnauthorized)
		return
	}
	modules.DB.Exec(`UPDATE sessions SET revoked = 1 WHERE session_id = ? AND refresh_token = ?`, sidCookie.Value, rtCookie.Value)
	// Clear session cookies
	ClearCookie(w, AuthInfo.JwtTokenName)
	ClearCookie(w, AuthInfo.SessionIDName)
	ClearCookie(w, AuthInfo.RefreshTokenName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
