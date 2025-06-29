package auth

import (
	"log"
	"net/http"
	"time"

	"social-network/app/Auth/jwt"
	database "social-network/app/modules"
	"social-network/server/logs"

	"github.com/google/uuid"
)

func CheckSession(r *http.Request, userID int, username string) (jwtToken, sessionID, refreshToken string, err error) {
	InvalidateSessions(userID)

	userAgent := r.Header.Get("User-Agent")
	ip := jwt.GetIP(r)

	// 1. Create DB session with refresh token
	sessionID, refreshToken, err = createSession(userID, userAgent, ip)
	if err != nil {
		return "", "", "", err
	}

	// 2. Generate JWT
	payload := jwt.CreateJwtPayload(AuthExpiration.JwtExpiration, userID, sessionID)
	jwtToken = jwt.Generate(payload)

	return jwtToken, sessionID, refreshToken, nil
}

func createSession(userID int, userAgent string, ip string) (string, string, error) {
	sessionID := uuid.New().String()
	refreshToken := jwt.GenerateToken(32) // 32 random bytes

	expiresAt := time.Now().Add(AuthExpiration.SessionExpiration)

	_, err := database.DB.Exec(`
		INSERT INTO sessions (user_id, session_id, refresh_token, expires_at, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, sessionID, refreshToken, expiresAt, ip, userAgent)
	if err != nil {
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

func CheckActiveSession(userID int) ([]string, error) {
	var sessions []string
	rows, err := database.DB.Query(`
        SELECT session_id 
        FROM sessions 
        WHERE user_id = ?
		ORDER BY expires_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var SessionID string
		if err := rows.Scan(&SessionID); err != nil {
			return nil, err
		}
		sessions = append(sessions, SessionID)
	}
	return sessions, nil
}

func InvalidateSessions(userID int) error {
	// Fetch active session ID
	activeSessions, _ := CheckActiveSession(userID)

	// Get all sessions associated with the user
	if len(activeSessions) > 0 {
		for _, sessionID := range activeSessions {
			err := invalidateSession(sessionID)
			if err != nil {
				log.Printf("Error invalidating session %s: %v", sessionID, err)
			}
		}
	}
	return nil
}

func invalidateSession(session_id string) error {
	_, err := database.DB.Exec(`
        DELETE FROM sessions 
        WHERE session_id = ?
    `, session_id)
	return err
}
