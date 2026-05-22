package jwtprovider_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/befragment/template-go/internal/lib/auth/jwtprovider"
	"github.com/befragment/template-go/internal/lib/auth/principal"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWTProviderGenerate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		tname string

		secret string
		ttl    time.Duration

		userID string
		role   string

		wantUserID string
		wantRole   string
	}{
		{
			tname:      "user",
			secret:     "secret",
			ttl:        5 * time.Minute,
			userID:     "user-1",
			role:       "user",
			wantUserID: "user-1",
			wantRole:   "user",
		},
		{
			tname:      "admin",
			secret:     "secret",
			ttl:        10 * time.Minute,
			userID:     "user-2",
			role:       "admin",
			wantUserID: "user-2",
			wantRole:   "admin",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.tname, func(t *testing.T) {
			t.Parallel()

			p := jwtprovider.NewJWTProvider(tc.secret, tc.ttl)
			before := time.Now().UTC()
			token, err := p.Generate(ctx, tc.userID, tc.role)
			if err != nil {
				t.Fatalf("Generate returned error: %v", err)
			}
			if token == "" {
				t.Fatalf("expected non-empty token")
			}

			after := time.Now().UTC()

			var claims jwtprovider.Claims
			parsedToken, err := jwt.ParseWithClaims(
				token,
				&claims,
				func(tok *jwt.Token) (any, error) {
					if tok.Method.Alg() != jwt.SigningMethodHS256.Alg() {
						return nil, jwtprovider.ErrInvalidToken
					}
					return []byte(tc.secret), nil
				},
			)
			if err != nil {
				t.Fatalf("ParseWithClaims returned error: %v", err)
			}
			if parsedToken == nil || !parsedToken.Valid {
				t.Fatalf("expected token to be valid")
			}

			if claims.UserID != tc.wantUserID || claims.Role != tc.wantRole {
				t.Fatalf("expected claims user_id=%q role=%q, got user_id=%q role=%q",
					tc.wantUserID, tc.wantRole, claims.UserID, claims.Role,
				)
			}
			if claims.IssuedAt == nil || claims.ExpiresAt == nil {
				t.Fatalf("expected issued_at and expires_at to be set")
			}

			iat := claims.IssuedAt.Time
			exp := claims.ExpiresAt.Time

			const skew = 2 * time.Second
			if iat.Before(before.Add(-skew)) || iat.After(after.Add(skew)) {
				t.Fatalf("unexpected issued_at: got %v, window [%v..%v]", iat, before, after)
			}

			expectedExp := iat.Add(tc.ttl)
			if absDuration(exp.Sub(expectedExp)) > skew {
				t.Fatalf("unexpected expires_at: got %v, expected around %v (ttl=%v)", exp, expectedExp, tc.ttl)
			}
		})
	}
}

func TestJWTProviderParse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const secret = "secret"
	const otherSecret = "other-secret"

	now := time.Now().UTC()
	exp := now.Add(2 * time.Minute)
	expiredClaims := jwtprovider.Claims{
		UserID: "user-1",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Minute)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
		},
	}

	makeTokenWithClaims := func(signedSecret string, method jwt.SigningMethod, claims jwt.Claims) string {
		tok := jwt.NewWithClaims(method, claims)
		s, err := tok.SignedString([]byte(signedSecret))
		if err != nil {
			t.Fatalf("SignedString error: %v", err)
		}
		return s
	}

	makeValidClaims := func(userID, role string) jwtprovider.Claims {
		return jwtprovider.Claims{
			UserID: userID,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(exp),
			},
		}
	}

	tests := []struct {
		tname          string
		providerSecret string

		token string

		wantErrIs error
		want      *principal.Principal
	}{
		{
			tname:          "valid token",
			providerSecret: secret,
			token: makeTokenWithClaims(
				secret,
				jwt.SigningMethodHS256,
				makeValidClaims("user-1", "user"),
			),
			want: &principal.Principal{UserID: "user-1", Role: "user"},
		},
		{
			tname:          "invalid signature: other secret",
			providerSecret: secret,
			token: makeTokenWithClaims(
				otherSecret,
				jwt.SigningMethodHS256,
				makeValidClaims("user-1", "user"),
			),
			wantErrIs: jwtprovider.ErrInvalidToken,
		},
		{
			tname:          "invalid signing method: HS512",
			providerSecret: secret,
			token: makeTokenWithClaims(
				secret,
				jwt.SigningMethodHS512,
				makeValidClaims("user-1", "user"),
			),
			wantErrIs: jwtprovider.ErrInvalidToken,
		},
		{
			tname:          "empty user_id",
			providerSecret: secret,
			token: makeTokenWithClaims(
				secret,
				jwt.SigningMethodHS256,
				makeValidClaims("", "user"),
			),
			wantErrIs: jwtprovider.ErrInvalidToken,
		},
		{
			tname:          "empty role",
			providerSecret: secret,
			token: makeTokenWithClaims(
				secret,
				jwt.SigningMethodHS256,
				makeValidClaims("user-1", ""),
			),
			wantErrIs: jwtprovider.ErrInvalidToken,
		},
		{
			tname:          "expired token",
			providerSecret: secret,
			token: makeTokenWithClaims(
				secret,
				jwt.SigningMethodHS256,
				expiredClaims,
			),
			wantErrIs: jwtprovider.ErrInvalidToken,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.tname, func(t *testing.T) {
			t.Parallel()

			p := jwtprovider.NewJWTProvider(tc.providerSecret, 1*time.Hour)
			pr, err := p.Parse(ctx, tc.token)

			if tc.wantErrIs != nil {
				if pr != nil {
					t.Fatalf("expected nil principal, got: %#v", pr)
				}
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, tc.wantErrIs) {
					t.Fatalf("expected error to match %v, got: %v", tc.wantErrIs, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if pr == nil {
				t.Fatalf("expected non-nil principal")
			}
			if tc.want == nil {
				t.Fatalf("test misconfigured: want must be set for success case")
			}
			assertPrincipalEqual(t, *tc.want, *pr)
		})
	}
}

func assertPrincipalEqual(t *testing.T, expected principal.Principal, got principal.Principal) {
	t.Helper()

	if got.UserID != expected.UserID || got.Role != expected.Role {
		t.Fatalf("expected principal %#v, got %#v", expected, got)
	}
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
