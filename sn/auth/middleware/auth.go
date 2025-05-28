package middleware

import (
	"context"
	"fmt"
	"net/http"

	"social-network/sn/auth"
	"social-network/sn/auth/jwt"
)

// Auth Middleware Enforcing that the user is authenticated
func MdlwLogged_IN(next http.HandlerFunc) http.HandlerFunc {
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
		ctx := context.WithValue(r.Context(), auth.UserContextKey, payload)
		next(w, r.WithContext(ctx))
	}
}

// Auth Middleware Enforcing that the user is NOT authenticated
func MdlwLogged_OUT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cookie, err := r.Cookie("jwt")
		if err == nil {
			_, err := jwt.JWTVerify(cookie.Value)
			if err == nil {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, `{"error": "user is already logged in"}`)
				return
			}
		}

		next(w, r)
	}
}

func Mdlw_router(next, signed_out http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and verify JWT
		w.Header().Set("Content-Type", "application/json")
		cookie, err := r.Cookie("jwt")
		if err != nil {
			signed_out(w, r)
			return
		}

		payload, err := jwt.JWTVerify(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"error": "invalid or expired token"}`)
			return
		}

		// Set user in context and proceed
		ctx := context.WithValue(r.Context(), auth.UserContextKey, payload)
		next(w, r.WithContext(ctx))
	}
}
