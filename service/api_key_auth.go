package service

import (
	"fmt"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/internal"
	"github.com/responsible-api/responsible-auth/resource/access"
	"github.com/responsible-api/responsible-auth/resource/user"
	"github.com/responsible-api/responsible-auth/tools"

	"github.com/golang-jwt/jwt/v5"
)

type APIKeyAuth struct {
	auth.AuthProvider
}

func NewApiKeyAuth() auth.AuthInterface {
	var provider auth.AuthInterface = &APIKeyAuth{}
	return provider
}

// SetOptions sets the options for the APIKeyAuth provider.
func (d *APIKeyAuth) SetOptions(options auth.AuthOptions) {
	Options = options
}

func (d *APIKeyAuth) Decode(APIKey string) (string, string, error) {
	unpackedUsername, unpackedPassword, err := validateAPIKey(APIKey)
	if err != nil {
		return "", "", err
	}
	// Return the decoded username and password
	return unpackedUsername, unpackedPassword, nil
}

func (a *APIKeyAuth) CreateAccessToken(userID string, APIKey string) (*access.RToken, error) {
	db, err := tools.NewDatabase()
	if err != nil {
		return nil, err
	}

	userRepo := user.NewRepository(db)
	user, err := userRepo.Read(userID, APIKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("User found:", user)

	token, err := internal.CreateAccessToken(Options)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *APIKeyAuth) CreateRefreshToken(userID string, hash string) (*access.RToken, error) {
	db, err := tools.NewDatabase()
	if err != nil {
		return nil, err
	}

	userRepo := user.NewRepository(db)
	user, err := userRepo.Read(userID, hash)
	if err != nil {
		return nil, err
	}

	refreshToken, err := internal.CreateRefreshToken(user.Name, Options)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (a *APIKeyAuth) GrantRefreshToken(refreshTokenString string) (*access.RToken, error) {
	refreshToken, err := internal.GrantRefreshToken(refreshTokenString, Options)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (a *APIKeyAuth) Validate(tokenString string) (*jwt.Token, error) {
	token, err := internal.Validate(tokenString, Options)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func validateAPIKey(APIKey string) (string, string, error) {
	// Implement your API key validation logic here
	// For example, you might want to decode the API key and extract the username and password
	// or perform any other necessary validation.

	// This is just a placeholder implementation.
	unpackedUsername := "exampleUser"
	unpackedPassword := "examplePassword"

	return unpackedUsername, unpackedPassword, nil
}
