package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/light-speak/lighthouse/env"
)

type claim struct {
	UserId int64
	jwt.RegisteredClaims
}

var key []byte

func init() {
	key = []byte(env.Getenv("JWT_KEY", "IWY@*3JUI#d309HhefzX2WpLtPKtD!hn"))
}
