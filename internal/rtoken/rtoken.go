package rtoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RToken struct {
	*jwt.Token
}

func NewToken(token *jwt.Token) *RToken {
	return &RToken{
		Token: token,
	}
}

func (r *RToken) GetToken() string {
	return r.Token.Raw
}

func (r *RToken) GetExpirationTime() (*jwt.NumericDate, error) {
	return r.Token.Claims.GetExpirationTime()
}

func (r *RToken) GetIssuedAt() (time.Time, error) {
	if claims, ok := r.Claims.(jwt.MapClaims); ok {
		if iat, ok := claims["iat"].(float64); ok {
			return time.Unix(int64(iat), 0), nil
		}
	}
	return time.Time{}, fmt.Errorf("issued at time not found in token")
}

func (r *RToken) GetNotBefore() (time.Time, error) {
	if claims, ok := r.Claims.(jwt.MapClaims); ok {
		if nbf, ok := claims["nbf"].(float64); ok {
			return time.Unix(int64(nbf), 0), nil
		}
	}
	return time.Time{}, fmt.Errorf("not before time not found in token")
}
