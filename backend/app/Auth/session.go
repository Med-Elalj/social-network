package auth

import (
	"net/http"
	"time"

	"social-network/app/Auth/jwt"
	database "social-network/app/modules"
	"social-network/server/logs"

	"github.com/google/uuid"
)

func CheckSession(r *http.Request, userID int, username string) (jwtToken, sessionID, refreshToken string, err error) {
	invalidateSessions(userID)

	userAgent := r.Header.Get("User-Agent")
	ip := GetIP(r)

	// 1. Create DB session with refresh token
	sessionID, refreshToken, err = createSession(userID, userAgent, ip)
	if err != nil {
		return "", "", "", err
	}

	// 2. Generate JWT
	payload := jwt.CreateJwtPayload(AuthInfo.JwtExpiration, userID, username, sessionID)
	jwtToken = jwt.Generate(payload)
	return jwtToken, sessionID, refreshToken, nil
}

func createSession(userID int, userAgent string, ip string) (string, string, error) {
	sessionID := uuid.NewString()
	refreshToken := uuid.NewString()
	expiresAt := time.Now().Add(AuthInfo.SessionExpiration)

	_, err := database.DB.Exec(`
		INSERT INTO sessions (user_id, session_id, refresh_token, expires_at, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, sessionID, refreshToken, expiresAt, ip, userAgent)
	if err != nil {
		logs.ErrorLog.Println("DB error in createSession:", err)
		return "", "", err
	}
	return sessionID, refreshToken, nil
}

func SessionExists(userID int, sessionID string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM sessions WHERE session_id = ? AND user_id = ?"
	err := database.DB.QueryRow(query, sessionID, userID).Scan(&count)
	if err != nil {
		logs.ErrorLog.Println("DB error in SessionExists:", err)
		return false, err
	}
	return count == 1, nil
}

func invalidateSessions(user_id int) error {
	_, err := database.DB.Exec(`
        DELETE FROM sessions 
        WHERE user_id = ?
    `, user_id)
	return err
}
