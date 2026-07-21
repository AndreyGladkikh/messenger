package http_server

import (
	"messenger/messenger/internal/platform/auth"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("User")
    
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// userID := "76c5ed7b-0a33-4ae0-bd07-9beff16908f2"
		ctx := auth.NewContextWithUserID(r.Context(), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}