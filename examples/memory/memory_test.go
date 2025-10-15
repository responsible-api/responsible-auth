package memory

import (
	"testing"

	"github.com/responsible-api/responsible-auth/storage"
	"github.com/responsible-api/responsible-auth/testutils"
)

func TestNewInMemoryStorage(t *testing.T) {
	memStorage := NewInMemoryStorage()

	if memStorage == nil {
		t.Errorf("NewInMemoryStorage() returned nil")
	}

	// Test that it implements UserStorage interface
	var _ storage.UserStorage = memStorage
}

func TestInMemoryStorage_FindUserByCredentials(t *testing.T) {
	memStorage := NewInMemoryStorage()

	tests := []struct {
		name        string
		username    string
		credentials string
		expectError bool
		expectUser  bool
	}{
		{
			name:        "valid email and credentials",
			username:    "test@example.com",
			credentials: "ipHEh|$==*#59@|ftT;IER^qgGG_sz!w",
			expectError: false,
			expectUser:  true,
		},
		{
			name:        "valid account ID and credentials",
			username:    "123456789",
			credentials: "ipHEh|$==*#59@|ftT;IER^qgGG_sz!w",
			expectError: false,
			expectUser:  true,
		},
		{
			name:        "invalid username",
			username:    "nonexistent@example.com",
			credentials: "ipHEh|$==*#59@|ftT;IER^qgGG_sz!w",
			expectError: true,
			expectUser:  false,
		},
		{
			name:        "invalid credentials",
			username:    "test@example.com",
			credentials: "wrong-password",
			expectError: true,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := memStorage.FindUserByCredentials(tt.username, tt.credentials)

			if tt.expectError && err == nil {
				t.Errorf("FindUserByCredentials() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("FindUserByCredentials() unexpected error = %v", err)
				return
			}

			if tt.expectUser && user == nil {
				t.Errorf("FindUserByCredentials() expected user but got nil")
				return
			}

			if !tt.expectUser && user != nil {
				t.Errorf("FindUserByCredentials() expected nil user but got %v", user)
				return
			}

			if tt.expectUser && user != nil {
				if user.Mail != "test@example.com" {
					t.Errorf("FindUserByCredentials() user.Mail = %v, want %v", user.Mail, "test@example.com")
				}

				if user.AccountID != 123456789 {
					t.Errorf("FindUserByCredentials() user.AccountID = %v, want %v", user.AccountID, 123456789)
				}
			}
		})
	}
}

func TestInMemoryStorage_FindUserByAPIKey(t *testing.T) {
	memStorage := NewInMemoryStorage()

	tests := []struct {
		name        string
		apiKey      string
		expectError bool
		expectUser  bool
	}{
		{
			name:        "valid api key",
			apiKey:      "api_key_12345",
			expectError: false,
			expectUser:  true,
		},
		{
			name:        "invalid api key",
			apiKey:      "invalid_api_key",
			expectError: true,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := memStorage.FindUserByAPIKey(tt.apiKey)

			if tt.expectError && err == nil {
				t.Errorf("FindUserByAPIKey() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("FindUserByAPIKey() unexpected error = %v", err)
				return
			}

			if tt.expectUser && user == nil {
				t.Errorf("FindUserByAPIKey() expected user but got nil")
				return
			}

			if !tt.expectUser && user != nil {
				t.Errorf("FindUserByAPIKey() expected nil user but got %v", user)
				return
			}

			if tt.expectUser && user != nil {
				if user.APIKey != "api_key_12345" {
					t.Errorf("FindUserByAPIKey() user.APIKey = %v, want %v", user.APIKey, "api_key_12345")
				}
			}
		})
	}
}

func TestInMemoryStorage_UpdateRefreshToken(t *testing.T) {
	memStorage := NewInMemoryStorage()

	tests := []struct {
		name         string
		userID       string
		refreshToken string
		expectError  bool
	}{
		{
			name:         "valid user",
			userID:       "testuser",
			refreshToken: "new_refresh_token",
			expectError:  false,
		},
		{
			name:         "invalid user",
			userID:       "nonexistent",
			refreshToken: "new_refresh_token",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := memStorage.UpdateRefreshToken(tt.userID, tt.refreshToken)

			if tt.expectError && err == nil {
				t.Errorf("UpdateRefreshToken() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("UpdateRefreshToken() unexpected error = %v", err)
				return
			}

			// If successful, verify the token was stored
			if !tt.expectError && tt.refreshToken != "" {
				user, err := memStorage.ValidateRefreshToken(tt.refreshToken)
				if err != nil {
					t.Errorf("UpdateRefreshToken() token not stored properly: %v", err)
				}

				if user == nil {
					t.Errorf("UpdateRefreshToken() stored token returned nil user")
				}
			}
		})
	}
}

func TestInMemoryStorage_ValidateRefreshToken(t *testing.T) {
	memStorage := NewInMemoryStorage()

	// First, add a refresh token
	err := memStorage.UpdateRefreshToken("testuser", "valid_refresh_token")
	if err != nil {
		t.Fatalf("Failed to set up test: %v", err)
	}

	tests := []struct {
		name         string
		refreshToken string
		expectError  bool
		expectUser   bool
	}{
		{
			name:         "valid refresh token",
			refreshToken: "valid_refresh_token",
			expectError:  false,
			expectUser:   true,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid_refresh_token",
			expectError:  true,
			expectUser:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := memStorage.ValidateRefreshToken(tt.refreshToken)

			if tt.expectError && err == nil {
				t.Errorf("ValidateRefreshToken() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("ValidateRefreshToken() unexpected error = %v", err)
				return
			}

			if tt.expectUser && user == nil {
				t.Errorf("ValidateRefreshToken() expected user but got nil")
				return
			}

			if !tt.expectUser && user != nil {
				t.Errorf("ValidateRefreshToken() expected nil user but got %v", user)
				return
			}
		})
	}
}

func TestInMemoryStorage_Interface(t *testing.T) {
	// Test that InMemoryStorage implements the storage.UserStorage interface
	var userStorage storage.UserStorage = NewInMemoryStorage()

	// Test all interface methods exist and can be called
	testUser := testutils.TestUser()

	// Test FindUserByCredentials
	_, err := userStorage.FindUserByCredentials("test@example.com", "ipHEh|$==*#59@|ftT;IER^qgGG_sz!w")
	if err != nil {
		t.Errorf("Interface method FindUserByCredentials failed: %v", err)
	}

	// Test FindUserByAPIKey
	_, err = userStorage.FindUserByAPIKey("api_key_12345")
	if err != nil {
		t.Errorf("Interface method FindUserByAPIKey failed: %v", err)
	}

	// Test UpdateRefreshToken
	err = userStorage.UpdateRefreshToken(testUser.Name, "test_refresh_token")
	if err != nil {
		t.Errorf("Interface method UpdateRefreshToken failed: %v", err)
	}

	// Test ValidateRefreshToken
	_, err = userStorage.ValidateRefreshToken("test_refresh_token")
	if err != nil {
		t.Errorf("Interface method ValidateRefreshToken failed: %v", err)
	}
}
