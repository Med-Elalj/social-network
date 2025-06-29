package middleware

import (
	"context"
	"fmt"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
)

// Auth Middleware Enforcing that the user is authenticated
// func Logged_IN(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		cookie, err := r.Cookie("JWT")
// 		if err != nil {
// 			w.WriteHeader(http.StatusForbidden)
// 			fmt.Fprint(w, `{"error": "jwt cookie not found"}`)
// 			return
// 		}

// 		payload, err := jwt.JWTVerify(cookie.Value)
// 		if err != nil || payload == nil {
// 			w.WriteHeader(http.StatusForbidden)
// 			fmt.Fprint(w, `{"error": "invalid or expired token"}`)
// 			return
// 		}

// 		// Attach payload (user info) to context
// 		ctx := context.WithValue(r.Context(), auth.UserContextKey, payload)
// 		next(w, r.WithContext(ctx))
// 	}
// }

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err1 := r.Cookie("JWT")
		ssidCookie, err2 := r.Cookie("ssid")
		if err1 != nil || err2 != nil {
			http.Error(w, "Unauthorized - missing cookies", http.StatusUnauthorized)
			return
		}

		payload, err := jwt.JWTVerify(jwtCookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized - invalid JWT", http.StatusUnauthorized)
			return
		}

		// Check if JWT's session ID matches cookie
		if payload.SessionID != ssidCookie.Value {
			http.Error(w, "Unauthorized - session ID mismatch", http.StatusUnauthorized)
			return
		}

		// Query DB for session details (including stored IP and User-Agent)
		session, err := auth.GetSessionByID(payload.SessionID)
		if err != nil {
			http.Error(w, "Unauthorized - session not found", http.StatusUnauthorized)
			return
		}
		fmt.Println("MiddleWare Session details:", session)
		// Optional: verify IP & User-Agent bind to session
		clientIP := auth.GetIP(r)
		if session.IP != clientIP {
			http.Error(w, "Unauthorized - IP mismatch", http.StatusUnauthorized)
			return
		}
		if session.UserAgent != r.UserAgent() {
			http.Error(w, "Unauthorized - User-Agent mismatch", http.StatusUnauthorized)
			return
		}

		// Put user info in context for handlers to use
		ctx := context.WithValue(r.Context(), auth.UserContextKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
