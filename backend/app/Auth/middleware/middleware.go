package middleware

import (
	"context"
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err1 := r.Cookie(auth.AuthInfo.JwtTokenName)
		ssidCookie, err2 := r.Cookie(auth.AuthInfo.SessionIDName)
		if err1 != nil || err2 != nil {
			auth.JsRespond(w, "Unauthorized - missing cookies", http.StatusUnauthorized)
			return
		}

		payload, err := jwt.JWTVerify(jwtCookie.Value)
		if err != nil {
			auth.JsRespond(w, "Unauthorized - invalid JWT", http.StatusUnauthorized)
			return
		}

		// Check if JWT's session ID matches cookie
		if payload.SessionID != ssidCookie.Value {
			auth.JsRespond(w, "Unauthorized - session ID mismatch", http.StatusUnauthorized)
			return
		}

		// Query DB for session details (including stored IP and User-Agent)
		session, err := auth.GetSessionByID(payload.SessionID)
		if err != nil {
			auth.ClearCookie(w, auth.AuthInfo.JwtTokenName)
			auth.JsRespond(w, "Unauthorized - session not found", http.StatusUnauthorized)
			return
		}
		// fmt.Println("MiddleWare Session details:", session)
		// verify IP & User-Agent bind to session
		clientIP := auth.GetIP(r)
		if session.IP != clientIP {
			auth.JsRespond(w, "Unauthorized - IP mismatch", http.StatusUnauthorized)
			auth.ClearCookie(w, auth.AuthInfo.JwtTokenName)
			return
		}
		if session.UserAgent != r.UserAgent() {
			auth.JsRespond(w, "Unauthorized - User-Agent mismatch", http.StatusUnauthorized)
			auth.ClearCookie(w, auth.AuthInfo.JwtTokenName)
			return
		}
		// Put user info in context for handlers to use
		ctx := context.WithValue(r.Context(), auth.UserContextKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow only specific origin (needed for credentials)
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		// Handle preflight (OPTIONS) request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Continue to actual handler
		next.ServeHTTP(w, r)
	})
}
