package middlewares

import "net/http"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if err == http.ErrNoCookie {
				// Cookie not found
				http.Error(w, "No session cookie found", http.StatusUnauthorized)
				return
			}
			// Other error
			http.Error(w, "Error reading cookie", http.StatusBadRequest)
			return
		}

		_ = cookie

		next.ServeHTTP(w, r)
	})
}
