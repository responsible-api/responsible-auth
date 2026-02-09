package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/concerns"

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

	if claims := token.Claims; claims != nil {
		if !validExpiry(claims) {
			return nil, fmt.Errorf("token expired")
		}

		if !validNotBefore(claims) {
			return nil, fmt.Errorf("token not valid yet")
		}
	}
	return token, nil
}

func validExpiry(claims jwt.Claims) bool {
	exp, err := claims.GetExpirationTime()
	return err == nil && exp != nil && !exp.Time.Before(time.Now())
}

func validNotBefore(claims jwt.Claims) bool {
	nbf, err := claims.GetNotBefore()
	return err == nil && nbf != nil && !nbf.Time.After(time.Now())
}

func ValidateRefreshToken(tokenString string, options auth.AuthOptions) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(options.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	if claims := token.Claims; claims != nil {
		if !validExpiry(claims) {
			return nil, fmt.Errorf("refresh token expired")
		}
	}
	return token, nil
}
