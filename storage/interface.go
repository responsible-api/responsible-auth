package storage

import "github.com/responsible-api/responsible-auth/resource/user"

// UserStorage defines the interface that external applications must implement
// to provide user data storage for the authentication library.
// This allows the library to be storage-agnostic.
type UserStorage interface {
	// FindUserByCredentials retrieves a user by username/email and validates their credentials
	// username can be either email or account_id
	// credentials is the hash/secret/password depending on auth method
	FindUserByCredentials(username, credentials string) (*user.User, error)

	// FindUserByAPIKey retrieves a user by their API key
	FindUserByAPIKey(apiKey string) (*user.User, error)

	// UpdateRefreshToken stores a refresh token for a user
	UpdateRefreshToken(userID string, refreshToken string) error

	// ValidateRefreshToken checks if a refresh token is valid for a user
	ValidateRefreshToken(refreshToken string) (*user.User, error)
}
