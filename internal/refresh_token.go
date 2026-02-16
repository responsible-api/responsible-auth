package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/resource/access"

	"github.com/golang-jwt/jwt/v5"
)

func CreateRefreshToken(username string, options auth.AuthOptions) (*access.RToken, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(options.RefreshTokenDuration).Unix(),
	})

	tokenString, err := refreshToken.SignedString([]byte(options.SecretKey))
	if err != nil {
		return nil, err
	}

	// Set the raw token string to the JWT token from the signed process
	refreshToken.Raw = tokenString
	// Return the refresh token string
	return access.NewToken(refreshToken), nil
}

func GrantRefreshToken(refreshTokenString string, options auth.AuthOptions) (*access.RToken, error) {
	// Parse and verify the requested refresh token to grant a new access token
	refreshToken, err := ValidateRefreshToken(refreshTokenString, options)
	if err != nil || !refreshToken.Valid {
		log.Println("Error parsing refresh token:", err)
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Generate a new access token if refresh token is valid
	newAccessToken, err := CreateAccessToken(options)
	if err != nil {
		return nil, err
	}
	return newAccessToken, nil
}
