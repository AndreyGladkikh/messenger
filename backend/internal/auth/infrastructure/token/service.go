package token

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	authDomain "messenger/messenger/internal/auth/domain"
	"messenger/messenger/internal/auth/domain/session"
	"messenger/messenger/internal/platform/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

var ErrInvalidToken = fmt.Errorf("%w: invalid token", authDomain.ErrUnauthorized)

type Service struct {
	jwtSignKey   crypto.PrivateKey
	jwtVerifyKey crypto.PublicKey
}

func NewService() (*Service, error) {
	projRoot := utils.ProjectRoot()

	signBytes, err := os.ReadFile(projRoot + "/keys/private.pem")
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseEdPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	verifyBytes, err := os.ReadFile(projRoot + "/keys/public.pem")
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
	sid string
}

func (s *Service) GenerateAccessToken(session *session.Session) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("EdDSA"))

	t.Claims = &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   session.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		sid: session.ID.String(),
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
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (any, error) {
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
