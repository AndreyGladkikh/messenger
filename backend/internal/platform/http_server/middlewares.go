package http_server

import (
	"fmt"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/messaging/infrastructure/auth"
	"net/http"

	"github.com/google/uuid"
)

func AuthMiddleware(tokenService *token.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := tokenService.ParseAccessTokenFromRequestAndGetClaims(r)
			if err != nil {
				NewResponse(nil, err).WriteTo(w)
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				err = fmt.Errorf("%w: %w", token.ErrInvalidToken, err)
				NewResponse(nil, err).WriteTo(w)
				return
			}

			// userID := r.Header.Get("User")
			// if userID == "" {
			// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// 	return
			// }
			// uID, err := uuid.Parse(userID)
			// if err != nil {
			// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// 	return
			// }
			// userID := "76c5ed7b-0a33-4ae0-bd07-9beff16908f2"

			ctx := auth.NewContextWithUserID(r.Context(), userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
