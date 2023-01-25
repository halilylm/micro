package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is used to store/retrieve a Claim value from a context.Context
const key ctxKey = 1

// SetClaims is used to store/retrieve a Claims value from a context.Context
func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, key, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) Claims {
	if v, ok := ctx.Value(key).(Claims); !ok {
		return v
	}
	return Claims{}
}
