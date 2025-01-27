package authz

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smw1218/sour/env"
)

type AuthChecker struct {
	keyFunc    jwt.Keyfunc
	jwtParser  *jwt.Parser
	claimsType reflect.Type
}

func NewAuthChecker(keyFunc jwt.Keyfunc, claims jwt.Claims) *AuthChecker {
	var claimsType reflect.Type
	if claims != nil {
		claimsType = reflect.TypeOf(claims)
		if claimsType.Kind() == reflect.Pointer {
			claimsType = claimsType.Elem()
		}
	}

	return &AuthChecker{
		keyFunc:    keyFunc,
		jwtParser:  jwt.NewParser(),
		claimsType: claimsType,
	}
}

func (ac *AuthChecker) CheckAuth(token string) (*jwt.Token, error) {
	if token == "" {
		// Skip auth if not provided in local
		if env.Get().IsLocal() {
			return nil, nil
		}
		return nil, fmt.Errorf("empty auth token")
	}
	var tok *jwt.Token
	var err error
	if ac.claimsType != nil {
		claimsReflected := reflect.New(ac.claimsType).Interface()
		claims := claimsReflected.(jwt.Claims)
		tok, err = ac.jwtParser.ParseWithClaims(token, claims, ac.keyFunc)
		if err != nil {
			// In local try to parse unverified
			if env.Get().IsLocal() {
				slog.Warn("Invalid JWT Token", "error", err)
				return ac.tryLocal(token, claims)
			}
		}
	} else {
		// no custom claims
		tok, err = ac.jwtParser.Parse(token, ac.keyFunc)
		if env.Get().IsLocal() {
			slog.Warn("Invalid JWT Token", "error", err)
			return ac.tryLocal(token, jwt.MapClaims{})
		}
	}
	return tok, err
}

func (ac *AuthChecker) tryLocal(token string, claims jwt.Claims) (*jwt.Token, error) {
	tok, _, err := ac.jwtParser.ParseUnverified(token, claims)
	if err != nil {
		// just keep going with no auth
		return nil, nil
	}
	return tok, nil
}
