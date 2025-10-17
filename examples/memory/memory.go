package memory

import (
	"errors"

	"github.com/responsible-api/responsible-auth/resource/user"
	"github.com/responsible-api/responsible-auth/storage"
)

// InMemoryStorage is a simple in-memory implementation of UserStorage
// This is useful for testing or applications that don't need persistent storage
type InMemoryStorage struct {
	users         map[string]*user.User // keyed by username/email
	apiKeys       map[string]*user.User // keyed by API key
	refreshTokens map[string]*user.User // keyed by refresh token
}

// NewInMemoryStorage creates a new in-memory storage with some sample data
func NewInMemoryStorage() storage.UserStorage {
	storage := &InMemoryStorage{
		users:         make(map[string]*user.User),
		apiKeys:       make(map[string]*user.User),
		refreshTokens: make(map[string]*user.User),
	}

	// Add sample users for demonstration
	sampleUser := &user.User{
		AccountID: 123456789,
		Name:      "test-user",
		Mail:      "test@example.com",
		Secret:    "ipHEh|$==*#59@|ftT;IER^qgGG_sz!w", // matches the decoded credentials
		APIKey:    "api_key_12345",
		Status:    1, // active
	}

	storage.users["test@example.com"] = sampleUser
	storage.users["test-user"] = sampleUser
	storage.users["123456789"] = sampleUser
	storage.apiKeys["api_key_12345"] = sampleUser

	return storage
}

// FindUserByCredentials retrieves a user by username/email and validates their credentials
func (m *InMemoryStorage) FindUserByCredentials(username, credentials string) (*user.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	if user.Secret != credentials {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// FindUserByAPIKey retrieves a user by their API key
func (m *InMemoryStorage) FindUserByAPIKey(apiKey string) (*user.User, error) {
	user, exists := m.apiKeys[apiKey]
	if !exists {
		return nil, errors.New("invalid API key")
	}
	return user, nil
}

// UpdateRefreshToken stores a refresh token for a user
func (m *InMemoryStorage) UpdateRefreshToken(userID string, refreshToken string) error {
	user, exists := m.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.Refresh = refreshToken
	m.refreshTokens[refreshToken] = user
	return nil
}

// ValidateRefreshToken checks if a refresh token is valid for a user
func (m *InMemoryStorage) ValidateRefreshToken(refreshToken string) (*user.User, error) {
	user, exists := m.refreshTokens[refreshToken]
	if !exists {
		return nil, errors.New("invalid refresh token")
	}

	return user, nil
}
