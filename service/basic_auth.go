package service

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/internal"
	"github.com/responsible-api/responsible-auth/resource/access"
	"github.com/responsible-api/responsible-auth/storage"

	"github.com/golang-jwt/jwt/v5"
)

var Options auth.AuthOptions

type BasicAuth struct {
	auth.AuthProvider
	storage storage.UserStorage
}

type AuthOptions struct {
	Options auth.AuthOptions
}

func NewBasicAuth() auth.AuthInterface {
	var provider auth.AuthInterface = &BasicAuth{}
	return provider
}

// SetOptions sets the options for the BasicAuth provider.
func (d *BasicAuth) SetOptions(options auth.AuthOptions) {
	Options = options
}

// SetStorage sets the storage implementation for the BasicAuth provider.
func (d *BasicAuth) SetStorage(storage storage.UserStorage) {
	d.storage = storage
}

func (d *BasicAuth) Decode(hash string) (string, string, error) {
	unpackedUsername, unpackedPassword, err := validateBasic(hash)
	if err != nil {
		return "", "", err
	}
	// Return the decoded username and password
	return unpackedUsername, unpackedPassword, nil
}

// Grant generates a token for the user with the given ID and password.
func (a *BasicAuth) CreateAccessToken(userID string, hash string) (*access.RToken, error) {
	_, err := a.storage.FindUserByCredentials(userID, hash)
	if err != nil {
		return nil, err
	}
	token, err := internal.CreateAccessToken(Options)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *BasicAuth) CreateRefreshToken(userID string, hash string) (*access.RToken, error) {
	user, err := a.storage.FindUserByCredentials(userID, hash)
	if err != nil {
		return nil, err
	}

	refreshToken, err := internal.CreateRefreshToken(user.Name, Options)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (a *BasicAuth) GrantRefreshToken(refreshTokenString string) (*access.RToken, error) {
	refreshToken, err := internal.GrantRefreshToken(refreshTokenString, Options)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (a *BasicAuth) Validate(tokenString string) (*jwt.Token, error) {
	token, err := internal.Validate(tokenString, Options)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// BasicAuth decodes a base64-encoded client credentials string and returns the username and password.
func validateBasic(encodedCredentials string) (string, string, error) {
	// Decode the base64-encoded string
	decoded, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return "", "", errors.New("invalid base64 encoding")
	}

	// Split the decoded string into username and password
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid credentials format")
	}
	username, password := parts[0], parts[1]

	if (username == "") || (password == "") {
		return "", "", errors.New("invalid credentials format")
	}

	// Return the decoded username and password
	return username, password, nil
}
