package testutils

import (
	"time"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/resource/user"
)

// TestAuthOptions returns a standard set of auth options for testing
func TestAuthOptions() auth.AuthOptions {
	return auth.AuthOptions{
		SecretKey:            "test-secret-key-32-characters!",
		TokenDuration:        1 * time.Hour,
		RefreshTokenDuration: 24 * time.Hour,
		TokenLeeway:          30 * time.Second,
		CookieDuration:       24 * time.Hour,
		Issuer:               "test-issuer",
		IssuedAt:             time.Now().Unix(),
		NotBefore:            time.Now().Unix(),
		Subject:              "test-subject",
		Scopes:               "read,write",
		Role:                 "user",
		CustomClaims: map[string]interface{}{
			"test_claim": "test_value",
			"number":     123,
		},
	}
}

// TestUser returns a test user for testing
func TestUser() *user.User {
	return &user.User{
		AccountID: 123456789,
		Name:      "testuser",
		Mail:      "test@example.com",
		Created:   uint64(time.Now().Unix()),
		Access:    uint64(time.Now().Unix()),
		Status:    1,
		Secret:    "test-password-hash",
		APIKey:    "test-api-key-12345",
		Refresh:   "",
	}
}

// TestUserWithRefresh returns a test user with a refresh token
func TestUserWithRefresh(refreshToken string) *user.User {
	u := TestUser()
	u.Refresh = refreshToken
	return u
}

// ValidBasicAuthCredentials returns a valid base64-encoded basic auth string
// Encodes "test@example.com:test-password-hash"
func ValidBasicAuthCredentials() string {
	return "dGVzdEBleGFtcGxlLmNvbTp0ZXN0LXBhc3N3b3JkLWhhc2g="
}

// InvalidBasicAuthCredentials returns various invalid basic auth strings for testing
func InvalidBasicAuthCredentials() []string {
	return []string{
		"invalid-base64!@#",
		"",         // empty string
		"dGVzdA==", // valid base64 but no colon separator
		"OnRlc3Q=", // starts with colon
		"dGVzdDo=", // ends with colon
	}
}

// MockStorage creates a simple mock storage for testing
type MockStorage struct {
	Users         map[string]*user.User
	APIKeys       map[string]*user.User
	RefreshTokens map[string]*user.User
	ShouldError   bool
	ErrorMessage  string
}

func NewMockStorage() *MockStorage {
	testUser := TestUser()
	return &MockStorage{
		Users: map[string]*user.User{
			"test@example.com": testUser,
			"123456789":        testUser,
		},
		APIKeys: map[string]*user.User{
			"test-api-key-12345": testUser,
		},
		RefreshTokens: make(map[string]*user.User),
		ShouldError:   false,
	}
}

func (m *MockStorage) FindUserByCredentials(username, credentials string) (*user.User, error) {
	if m.ShouldError {
		return nil, &TestError{Message: m.ErrorMessage}
	}

	if user, exists := m.Users[username]; exists && user.Secret == credentials {
		return user, nil
	}
	return nil, &TestError{Message: "user not found"}
}

func (m *MockStorage) FindUserByAPIKey(apiKey string) (*user.User, error) {
	if m.ShouldError {
		return nil, &TestError{Message: m.ErrorMessage}
	}

	if user, exists := m.APIKeys[apiKey]; exists {
		return user, nil
	}
	return nil, &TestError{Message: "api key not found"}
}

func (m *MockStorage) UpdateRefreshToken(userID string, refreshToken string) error {
	if m.ShouldError {
		return &TestError{Message: m.ErrorMessage}
	}

	// Find user by ID in Users map
	for _, user := range m.Users {
		if user.Name == userID {
			user.Refresh = refreshToken
			m.RefreshTokens[refreshToken] = user
			return nil
		}
	}
	return &TestError{Message: "user not found for refresh token update"}
}

func (m *MockStorage) ValidateRefreshToken(refreshToken string) (*user.User, error) {
	if m.ShouldError {
		return nil, &TestError{Message: m.ErrorMessage}
	}

	if user, exists := m.RefreshTokens[refreshToken]; exists {
		return user, nil
	}
	return nil, &TestError{Message: "refresh token not found"}
}

// SetError configures the mock to return errors
func (m *MockStorage) SetError(shouldError bool, message string) {
	m.ShouldError = shouldError
	m.ErrorMessage = message
}

// TestError is a simple error type for testing
type TestError struct {
	Message string
}

func (e *TestError) Error() string {
	return e.Message
}
