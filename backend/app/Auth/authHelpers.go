package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"social-network/app/modules"
	"social-network/server/logs"
)

type contextKey string

const UserContextKey contextKey = "user"

type authExpiration struct {
	JwtExpiration          time.Duration
	RefreshTokenExpiration time.Duration
	SessionExpiration      time.Duration
}

var AuthExpiration = authExpiration{
	JwtExpiration:          time.Duration(15 * time.Minute),   // 15 minutes
	SessionExpiration:      time.Duration(7 * 24 * time.Hour), // 7 days
	RefreshTokenExpiration: time.Duration(7 * 24 * time.Hour), // 7 days
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
