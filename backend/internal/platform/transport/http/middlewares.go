package http

import (
	"context"
	"fmt"
	authDomain "messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/infrastructure/token"
	"messenger/messenger/internal/platform/auth"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"messenger/messenger/internal/platform/transport/http/headers"
	"messenger/messenger/internal/shared/application/apperr"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"

	"uuid"
)

func AuthMiddleware(tokenService *token.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			switch {
			case isHeaderPresentInRequest(r, headers.Authorization),
				isCookiePresentInRequest(r, cookies.AccessToken):
				claims, err := tokenService.ParseAccessTokenFromRequestAndGetClaims(r)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				userID, err := uuid.Parse(claims.Subject)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				ctx = auth.NewContextWithAuth(ctx, auth.Auth{
					Method: auth.MethodBearerToken,
					Principal: auth.Principal{
						Subject: auth.Subject{
							ID:   userID,
							Type: auth.SubjectTypeUser,
						},
						Scope: nil,
					},
					SessionID: claims.SID,
				})

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			case isHeaderPresentInRequest(r, headers.XApiKey):
				//todo api key auth
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isHeaderPresentInRequest(r *http.Request, header string) bool {
	h := r.Header.Get(header)
	return h != ""
}

func isCookiePresentInRequest(r *http.Request, cookieName string) bool {
	_, err := r.Cookie(cookieName)
	return err == nil
}

func FirewallMiddleware(oas *openapi3.T, errHandler *ErrorHandler) func(next http.Handler) http.Handler {
	validator := oapimiddleware.OapiRequestValidatorWithOptions(oas,
		&oapimiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
					a, err := auth.AuthFromContext(ctx)
					if err != nil {
						err = fmt.Errorf("%w: %w", authDomain.ErrUnauthorized, err)
						return apperr.Map(err)
					}

					if input.SecuritySchemeName != string(a.Method) {
						return fmt.Errorf("security scheme %s != '%s'", input.SecuritySchemeName, string(a.Method))
					}

					err = CheckScopes(input.Scopes, a.Principal.Scope)
					if err != nil {
						return err
					}

					return nil
				},
			},
			ErrorHandlerWithOpts: func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request, opts oapimiddleware.ErrorHandlerOpts) {
				errHandler.Handle(w, r, err)
			},
		})

	return validator
}

func CheckScopes(expectedScopes, tokenScopes []string) error {
	scopesMap := make(map[string]bool, len(tokenScopes))
	for _, c := range tokenScopes {
		scopesMap[c] = true
	}

	for _, e := range expectedScopes {
		if !scopesMap[e] {
			// todo return access denied err
			return fmt.Errorf("access denied")
		}
	}
	return nil
}
