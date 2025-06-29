package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social-network/app/modules"
	"social-network/server/logs"
)

type contextKey string

const UserContextKey contextKey = "user"

type authInfo struct {
	JwtExpiration          time.Duration
	RefreshTokenExpiration time.Duration
	SessionExpiration      time.Duration
	JwtTokenName           string
	SessionIDName          string
	RefreshTokenName       string
}

var AuthInfo = authInfo{
	JwtExpiration:          time.Duration(15 * time.Minute),   // 15 minutes
	SessionExpiration:      time.Duration(7 * 24 * time.Hour), // 7 days
	RefreshTokenExpiration: time.Duration(7 * 24 * time.Hour), // 7 days
	JwtTokenName:           "JWT",
	SessionIDName:          "session_id",
	RefreshTokenName:       "refresh_token",
}

type Session struct {
	UserID       int
	SessionID    string
	RefreshToken string
	IP           string
	UserAgent    string
	Revoked      bool

	ExpiresAt time.Time
}

type responseMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   maxAge,
	})
}

func ClearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func GetElemVal(selectedElem, from, where string) any {
	var res any

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", selectedElem, from, where)

	err := modules.DB.QueryRow(query).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			res = ""
		} else {
			logs.ErrorLog.Println("Database error:", err)
		}
	}

	return res
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

func JsRespond(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(responseMessage{
		Message: message,
		Code:    code,
	})
}

func GetSessionByID(sessionID string) (Session, error) {
	var session Session
	modules.DB.QueryRow(`
	SELECT session_id, user_id, refresh_token, ip_address, user_agent, expires_at
	FROM sessions
	WHERE session_id = ?
	`, sessionID).Scan(
		&session.SessionID,
		&session.UserID,
		&session.RefreshToken,
		&session.IP,
		&session.UserAgent,
		&session.ExpiresAt,
	)
	if session.SessionID == "" {
		return session, fmt.Errorf("session not found")
	}
	return session, nil
}
