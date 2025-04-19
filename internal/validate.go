package internal

import (
	"fmt"
	"time"

	"github.com/vince-scarpa/responsible-api-go/auth"
	"github.com/vince-scarpa/responsible-api-go/concerns"

	"github.com/golang-jwt/jwt/v5"
)

func Validate(tokenString string, options auth.AuthOptions) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &concerns.ClaimsGeneric{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return token, nil
		}
		return []byte(options.SecretKey), nil
	}, jwt.WithLeeway(options.TokenLeeway))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(*concerns.ClaimsGeneric); ok && token.Valid {
		if !validExpiry(claims) {
			return nil, fmt.Errorf("token expired")
		}

		if !validNotBefore(claims) {
			return nil, fmt.Errorf("token not valid yet")
		}
	}
	return token, nil
}

func validExpiry(claims *concerns.ClaimsGeneric) bool {
	return !(claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()))
}

func validNotBefore(claims *concerns.ClaimsGeneric) bool {
	return !(claims.NotBefore == nil || claims.NotBefore.Time.After(time.Now()))
}
