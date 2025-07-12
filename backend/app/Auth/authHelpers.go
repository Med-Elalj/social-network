package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"social-network/app/modules"
	"social-network/server/logs"

	"golang.org/x/crypto/bcrypt"
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
type Register struct {
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
	Fname     string `json:"fname"`
	Lname     string `json:"lname"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	About     string `json:"about"`
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
		SameSite: http.SameSiteNoneMode,
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

func GetElemVal[T any](selectedElem, from, whereClause string, args ...any) (T, error) {
	var res T

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", selectedElem, from, whereClause)

	err := modules.DB.QueryRow(query, args...).Scan(&res)
	if err != nil {
		var zero T
		if err == sql.ErrNoRows {
			return zero, nil // no result, return empty value
		}
		logs.ErrorLog.Println("Database error:", err)
		return zero, err
	}

	return res, nil
}

func EntryExists(elem, value, from string, checkLower bool) (int, bool) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", from, elem)
	if checkLower {
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE LOWER(%s) = LOWER(?)", from, elem)
	}

	err := modules.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		logs.ErrorLog.Println("Database error:", err)
		return -1, false
	}

	return count, count > 0
}

func GetIP(r *http.Request) string {
	// todo: if hosted, remove comment
	// if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
	// 	return strings.TrimSpace(strings.Split(ip, ",")[0])
	// }
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
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

func ChangePassword(password string, userID int64) {
	query := `UPDATE user SET password_hash = ? WHERE id = ?`
	_, err := modules.DB.Exec(query, HashPassword(password), userID)
	if err != nil {
		logs.ErrorLog.Fatalln("Failed to change password:", err)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logs.ErrorLog.Fatalln(err.Error())
		return ""
	}
	return string(bytes)
}

func CheckPassword(password string, userID int) bool {
	var hashedPassword string
	query := `SELECT password_hash FROM user WHERE id = ?`
	err := modules.DB.QueryRow(query, userID).Scan(&hashedPassword)
	if err != nil {
		logs.ErrorLog.Fatalln(err.Error())
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func JsRespond(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(responseMessage{
		Message: message,
		Code:    code,
	})
}
