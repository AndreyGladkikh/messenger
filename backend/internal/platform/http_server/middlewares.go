package http_server

import (
	"fmt"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/platform/http_server/auth"
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
			sessionID, err := uuid.Parse(claims.SID)
			if err != nil {
				err = fmt.Errorf("%w: %w", token.ErrInvalidToken, err)
				NewResponse(nil, err).WriteTo(w)
				return
			}

			ctx := auth.NewContextWithUserID(r.Context(), userID)
			ctx = auth.NewContextWithSessionID(ctx, sessionID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
