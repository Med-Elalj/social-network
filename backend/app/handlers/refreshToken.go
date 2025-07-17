package handlers

import (
	"net/http"
	"strconv"
	"time"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/modules"

	"github.com/google/uuid"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	sidCookie, err1 := r.Cookie(auth.AuthInfo.SessionIDName)
	rtCookie, err2 := r.Cookie(auth.AuthInfo.RefreshTokenName)
	if err1 != nil || err2 != nil {
		auth.JsResponse(w, "Missing cookies", http.StatusUnauthorized)
		return
	}
	payload, ok := r.Context().Value(auth.UserContextKey).(*jwt.JwtPayload)
	if !ok {
		auth.JsResponse(w, "Unauthorized - invalid user", http.StatusUnauthorized)
		return
	}
	if payload.SessionID != sidCookie.Value {
		auth.JsResponse(w, "Unauthorized - session ID mismatch", http.StatusUnauthorized)
		return
	}
	// Validate session ID and refresh token
	session, err := auth.GetSessionByID(payload.SessionID)
	// Check if session exists and is valid
	if err != nil {
		auth.JsResponse(w, "Session not found", http.StatusUnauthorized)
		return
	}
	if session.RefreshToken != rtCookie.Value {
		auth.JsResponse(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	if time.Now().After(session.ExpiresAt) {
		auth.JsResponse(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}
	if auth.GetIP(r) != session.IP {
		auth.JsResponse(w, "IP address mismatch", http.StatusUnauthorized)
		return
	}
	if r.Header.Get("User-Agent") != session.UserAgent {
		auth.JsResponse(w, "User-Agent mismatch", http.StatusUnauthorized)
		return
	}

	// Rotate refresh token
	newRefreshToken := uuid.NewString()
	newExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = modules.DB.Exec(`
		UPDATE sessions SET refresh_token = ?, expires_at = ? WHERE session_id = ?
	`, newRefreshToken, newExpiresAt, sidCookie.Value)
	if err != nil {
		auth.JsResponse(w, "Failed to update session", http.StatusInternalServerError)
		return
	}

	// Generate new access token (JWT)
	username, _ := auth.GetElemVal[string]("display_name", "profile", "id="+strconv.Itoa(session.UserID))
	jwtToken := jwt.Generate(jwt.CreateJwtPayload(auth.AuthInfo.JwtExpiration, session.UserID, username, session.SessionID))

	auth.SetCookie(w, auth.AuthInfo.JwtTokenName, jwtToken, int(auth.AuthInfo.JwtExpiration.Seconds())) // 15 min
	auth.SetCookie(w, auth.AuthInfo.RefreshTokenName, newRefreshToken,
		int(auth.AuthInfo.RefreshTokenExpiration.Seconds()-auth.AuthInfo.JwtExpiration.Seconds())) // 7 days

	auth.JsResponse(w, "Refresh token updated", http.StatusOK)
}
