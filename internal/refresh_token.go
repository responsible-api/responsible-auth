package internal

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vince-scarpa/responsible-api-go/auth"
	"github.com/vince-scarpa/responsible-api-go/internal/rtoken"

	"github.com/golang-jwt/jwt/v5"
)

func CreateRefreshToken(username string, options auth.AuthOptions) (*rtoken.RToken, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(options.RefreshTokenDuration).Unix(),
	})

	_, err := refreshToken.SignedString([]byte(options.SecretKey))
	if err != nil {
		return nil, err
	}

	// Return the refresh token string
	return rtoken.NewToken(refreshToken), nil
}

func GrantRefreshToken(refreshTokenString string, options auth.AuthOptions) (*rtoken.RToken, error) {
	// Parse and verify the requested refresh token to grant a new access token
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(options.SecretKey), nil
	})

	if err != nil || !refreshToken.Valid {
		log.Println("Error parsing refresh token:", err)
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Generate a new access token if refresh token is valid
	if _, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		newAccessToken, err := CreateAccessToken(options)
		if err != nil {
			return nil, err
		}
		return newAccessToken, nil
	}

	// If the refresh token is not valid, return an error
	return nil, fmt.Errorf("invalid refresh token")
}
