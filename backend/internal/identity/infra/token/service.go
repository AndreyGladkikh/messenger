package token

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"messenger/messenger/internal/identity/app/token"
	"messenger/messenger/internal/identity/domain/model/session"
	"messenger/messenger/internal/platform/config"
	"messenger/messenger/internal/platform/transport/http/cookies"
	"messenger/messenger/internal/shared/app/apperr"
	"net/http"
	"os"
	"time"

	"uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

var ErrInvalidToken = apperr.NewUnauthenticated("Invalid token")

type Service struct {
	jwtSignKey   crypto.PrivateKey
	jwtVerifyKey crypto.PublicKey
}

func NewService(
	cfg config.Config,
) (*Service, error) {
	signBytes, err := os.ReadFile(cfg.Auth.Jwt.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseEdPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	verifyBytes, err := os.ReadFile(cfg.Auth.Jwt.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseEdPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	return &Service{
		jwtSignKey:   signKey,
		jwtVerifyKey: verifyKey,
	}, nil
}

type Claims struct {
	jwt.RegisteredClaims
	SID    uuid.UUID
	Scopes []string
}

func (s *Service) GenerateAccessToken(session *session.Session) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("EdDSA"))

	t.Claims = &Claims{
		Subject:   session.UserID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(token.AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		SID:       session.ID(),
	}

	return t.SignedString(s.jwtSignKey)
}

func (s *Service) GenerateRefreshToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (s *Service) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *Service) ParseAccessTokenFromRequestAndGetClaims(r *http.Request) (*Claims, error) {
	extractor := request.MultiExtractor{
		request.OAuth2Extractor,
		TokenCookieExtractor{},
	}
	token, err := request.ParseFromRequest(r, extractor, func(token *jwt.Token) (any, error) {
		return s.jwtVerifyKey, nil
	}, request.WithClaims(&Claims{}))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	return claims, nil
}

type TokenCookieExtractor struct{}

func (c TokenCookieExtractor) ExtractToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookies.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return "", request.ErrNoTokenInRequest
		default:
			return "", err
		}
	}
	return cookie.Value, nil
}
