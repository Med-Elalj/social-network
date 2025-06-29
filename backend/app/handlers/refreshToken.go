package handlers

import (
	"net/http"
	"time"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/modules"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	sidCookie, err1 := r.Cookie("session_id")
	rtCookie, err2 := r.Cookie("refresh_token")
	if err1 != nil || err2 != nil {
		http.Error(w, "Missing cookies", http.StatusUnauthorized)
		return
	}

	var session jwt.Session
	err := modules.DB.QueryRow(`
	SELECT session_id, user_id, refresh_token, ip_address, user_agent, expires_at
	FROM sessions
	WHERE session_id = ?
`, sidCookie.Value).Scan(
		&session.SessionID,
		&session.UserID,
		&session.RefreshToken,
		&session.IP,
		&session.UserAgent,
		&session.ExpiresAt,
	)
	if err != nil {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}
	if session.Revoked {
		http.Error(w, "Session revoked", http.StatusUnauthorized)
		return
	}
	if session.RefreshToken != rtCookie.Value {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	if time.Now().After(session.ExpiresAt) {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}
	if jwt.GetIP(r) != session.IP {
		http.Error(w, "IP address mismatch", http.StatusUnauthorized)
		return
	}
	if r.Header.Get("User-Agent") != session.UserAgent {
		http.Error(w, "User-Agent mismatch", http.StatusUnauthorized)
		return
	}

	// Rotate refresh token
	newRefreshToken := jwt.GenerateToken(32)
	newExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = modules.DB.Exec(`
		UPDATE sessions SET refresh_token = ?, expires_at = ? WHERE session_id = ?
	`, newRefreshToken, newExpiresAt, sidCookie.Value)
	if err != nil {
		http.Error(w, "Failed to update session", http.StatusInternalServerError)
		return
	}

	// Generate new access token (JWT)
	accessToken := jwt.Generate(jwt.CreateJwtPayload(auth.AuthExpiration.JwtExpiration, session.UserID, ""))
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	auth.SetCookie(w, "access_token", accessToken, int(auth.AuthExpiration.JwtExpiration.Seconds()))               // 15 min
	auth.SetCookie(w, "refresh_token", newRefreshToken, int(auth.AuthExpiration.RefreshTokenExpiration.Seconds())) // 7 days

	w.Write([]byte("Access token refreshed"))
}
