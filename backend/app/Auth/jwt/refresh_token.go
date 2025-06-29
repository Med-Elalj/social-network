package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

type Session struct {
	UserID       int
	SessionID    string
	RefreshToken string
	IP           string
	UserAgent    string
	Revoked      bool

	ExpiresAt time.Time
}

/*
// ValidateSession checks if User-Agent matches on refresh

	func ValidateSession(r *http.Request, session *Session) bool {
		requestUserAgent := r.Header.Get("User-Agent")
		if requestUserAgent != session.UserAgent {
			return false // User-Agent mismatch, possibly stolen token
		}
		if time.Now().After(session.ExpiresAt) {
			return false // Session expired
		}
		return true
	}

call this elsewhere):
Make sure to save the IP when storing a session:

ip := r.RemoteAddr

	Sessions[sessionID] = Session{
		UserID: userID,
		RefreshToken: refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IP: ip, // bind IP
	}
*/

func GenerateToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GetIP(r *http.Request) string {
	// Use headers if behind proxy
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
