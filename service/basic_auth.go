package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	auth "github.com/vince-scarpa/responsible-api-go"
	"github.com/vince-scarpa/responsible-api-go/internal"
	"github.com/vince-scarpa/responsible-api-go/resource/user"
	"github.com/vince-scarpa/responsible-api-go/tools"

	"github.com/golang-jwt/jwt/v5"
)

var Options auth.AuthOptions

type BasicAuth struct {
	auth.AuthProvider
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

func (d *BasicAuth) Decode(hash string) (string, string, error) {
	unpackedUsername, unpackedPassword, err := validateBasic(hash)
	if err != nil {
		return "", "", err
	}
	// Return the decoded username and password
	return unpackedUsername, unpackedPassword, nil
}

// Grant generates a token for the user with the given ID and password.
func (a *BasicAuth) Grant(userID string, hash string) (string, error) {
	db, err := tools.NewDatabase()
	if err != nil {
		return "", err
	}

	userRepo := user.NewRepository(db)
	user, err := userRepo.Read(userID)
	if err != nil {
		return "", err
	}
	fmt.Println("User found:", user)

	tokenString, err := internal.Grant(userID, hash, Options)
	if err != nil {
		return "", err
	}
	return tokenString, nil
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
