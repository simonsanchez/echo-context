package token

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	UserIDCtxKey = "user_id"
	ScopesCtxKey = "scopes"
)

var (
	defaultExpiry = time.Now().Add(15 * time.Minute)
	defaultScopes = []string{}
)

type Token struct {
	internal jwt.Token
}

func Default() (Token, error) {
	b := jwt.NewBuilder().
		Expiration(defaultExpiry).
		Claim(ScopesCtxKey, defaultScopes)

	token, err := b.Build()
	if err != nil {
		return Token{}, fmt.Errorf("unable to build jwt token: %w", err)
	}

	return Token{internal: token}, nil
}

func (t Token) Scopes() ([]string, error) {
	scopes, exists := t.internal.Get(ScopesCtxKey)
	if !exists {
		return nil, fmt.Errorf("token does not contain %q", ScopesCtxKey)
	}

	return scopes.([]string), nil
}
