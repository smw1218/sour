package authz

import (
	"encoding/json"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

func NewJWKSAuthChecker(claims jwt.Claims, issuers ...string) (*AuthChecker, error) {
	kf, err := keyfunc.NewDefault(issuers)
	if err != nil {
		return nil, fmt.Errorf("failed creating jwks keyfunc: %w", err)
	}
	return NewAuthChecker(kf.Keyfunc, claims), nil
}

func NewStaticJWKSChecker(claims jwt.Claims, raw json.RawMessage) (*AuthChecker, error) {
	kf, err := keyfunc.NewJWKSetJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("failed creating jwks keyfunc: %w", err)
	}
	return NewAuthChecker(kf.Keyfunc, claims), nil
}
