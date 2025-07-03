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
			auth.JsRespond(w, "Unauthorized - session not found", http.StatusUnauthorized)
			return
		}
		// fmt.Println("MiddleWare Session details:", session)
		// verify IP & User-Agent bind to session
		clientIP := auth.GetIP(r)
		if session.IP != clientIP {
			auth.JsRespond(w, "Unauthorized - IP mismatch", http.StatusUnauthorized)
			return
		}
		if session.UserAgent != r.UserAgent() {
			auth.JsRespond(w, "Unauthorized - User-Agent mismatch", http.StatusUnauthorized)
			return
		}
		// Put user info in context for handlers to use
		ctx := context.WithValue(r.Context(), auth.UserContextKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
