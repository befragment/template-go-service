package jwtprovider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-befragment/pkg/auth/principal"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTProvider struct {
	secret []byte
	ttl    time.Duration
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

func NewJWTProvider(secret string, ttl time.Duration) *JWTProvider {
	return &JWTProvider{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *JWTProvider) Generate(_ context.Context, userID, role string) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return signedToken, nil
}

func (s *JWTProvider) Parse(_ context.Context, tokenString string) (*principal.Principal, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}

			return s.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.UserID == "" || claims.Role == "" {
		return nil, ErrInvalidToken
	}

	return &principal.Principal{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
