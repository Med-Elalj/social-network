package middleware

// // Auth Middleware Enforcing that the user is authenticated
// func Logged_IN(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Extract and verify JWT
// 		w.Header().Set("Content-Type", "application/json")
// 		cookie, err := r.Cookie("JWT")
// 		if err != nil {
// 			w.WriteHeader(http.StatusForbidden)
// 			fmt.Fprint(w, `{"error": "jwt cookie not found"}`)
// 			return
// 		}

// 		payload, err := jwt.JWTVerify(cookie.Value)
// 		if err != nil {
// 			w.WriteHeader(http.StatusForbidden)
// 			fmt.Fprint(w, `{"error": "invalid or expired token"}`)
// 			return
// 		}

// 		// Set user in context and proceed
// 		ctx := context.WithValue(r.Context(), auth.UserContextKey, payload)
// 		next(w, r.WithContext(ctx))
// 	}
// }

// // Auth Middleware Enforcing that the user is NOT authenticated
// func Logged_OUT(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		cookie, err := r.Cookie("jwt")
// 		if err == nil {
// 			_, err := jwt.JWTVerify(cookie.Value)
// 			if err == nil {
// 				w.WriteHeader(http.StatusForbidden)
// 				fmt.Fprint(w, `{"error": "user is already logged in"}`)
// 				return
// 			}
// 		}

// 		next(w, r)
// 	}
// }
