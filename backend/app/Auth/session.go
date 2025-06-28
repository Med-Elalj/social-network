package auth

import (
	"log"
	"time"

	"social-network/app/Auth/jwt"
	database "social-network/app/modules"
	"social-network/server/logs"

	"github.com/gofrs/uuid"
)

// Todo: send errors in js notifications
// todo: check for session id

func CheckSession(userID int, username string) (string, string, error) {
	InvalidateSessions(userID)

	// Create a new session if there's no active session
	sessionID, err := createSession(userID)
	if err != nil {
		return "", "", err
	}

	payload := jwt.CreateJwtPayload(userID, username)
	jwtToken := jwt.Generate(payload)

	return jwtToken, sessionID, nil
}

func createSession(userID int) (string, error) {
	// Generate unique session ID (UUID or random string)
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	// Set session expiration time from int64
	expiresAt := time.Now().Add(time.Duration(jwt.Time_to_Expire))

	// Insert the session into the database
	_, err = database.DB.Exec(`
        INSERT INTO sessions (user_id, session_id , expires_at) 
        VALUES (?, ?, ?)
    `, userID, sessionID.String(), expiresAt)
	if err != nil {
		return "", err
	}

	return sessionID.String(), nil
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
