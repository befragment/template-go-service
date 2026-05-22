package jwtprovider

import (
	"context"

	"github.com/befragment/template-go/internal/lib/auth/principal"
)

type contextKey string

const principalKey contextKey = "principal"

func WithPrincipal(ctx context.Context, principal *principal.Principal) context.Context {
	return context.WithValue(ctx, principalKey, principal)
}

func PrincipalFromContext(ctx context.Context) (*principal.Principal, bool) {
	principal, ok := ctx.Value(principalKey).(*principal.Principal)
	return principal, ok
}
