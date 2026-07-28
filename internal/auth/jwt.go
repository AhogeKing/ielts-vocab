package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const ClaimsContextKey = "auth_claims"

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewTokenManagerFromEnv() (*TokenManager, error) {
	ttl, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		return nil, fmt.Errorf("parse JWT_TTL: %w", err)
	}

	return NewTokenManager(
		os.Getenv("JWT_SECRET"),
		os.Getenv("JWT_ISSUER"),
		ttl,
	)
}

func NewTokenManager(secret, issuer string, ttl time.Duration) (*TokenManager, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must contain at least 32 characters")
	}
	if strings.TrimSpace(issuer) == "" {
		return nil, errors.New("JWT_ISSUER must be set")
	}
	if ttl <= 0 {
		return nil, errors.New("JWT_TTL must be positive")
	}

	return &TokenManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}, nil
}

func (m *TokenManager) Generate(userID int64, username string) (AccessToken, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return AccessToken{}, fmt.Errorf("sign JWT: %w", err)
	}

	return AccessToken{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(m.ttl.Seconds()),
	}, nil
}

func (m *TokenManager) Parse(rawToken string) (Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		rawToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method %q", token.Header["alg"])
			}
			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, errors.New("invalid JWT")
	}

	return *claims, nil
}
