package auth

import (
	"context"
	"fmt"
	"net/http"

	"social-network/sn/auth/jwt"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and verify JWT
		w.Header().Set("Content-Type", "application/json")
		cookie, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"error": "jwt cookie not found"}`)
			return
		}

		payload, err := jwt.JWTVerify(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"error": "invalid or expired token"}`)
			return
		}

		// Set user in context and proceed
		ctx := context.WithValue(r.Context(), UserContextKey, payload)
		next(w, r.WithContext(ctx))
	}
}
