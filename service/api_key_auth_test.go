package service

import (
	"testing"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/testutils"
)

func TestNewApiKeyAuth(t *testing.T) {
	provider := NewApiKeyAuth()

	if provider == nil {
		t.Errorf("NewApiKeyAuth() returned nil")
	}

	// Test that it implements AuthInterface
	var _ auth.AuthInterface = provider
}

func TestAPIKeyAuth_SetOptions(t *testing.T) {
	provider := NewApiKeyAuth()
	options := testutils.TestAuthOptions()

	provider.SetOptions(options)

	// Verify options were set by checking the global Options variable
	if Options.SecretKey != options.SecretKey {
		t.Errorf("SetOptions() SecretKey = %v, want %v", Options.SecretKey, options.SecretKey)
	}

	if Options.TokenDuration != options.TokenDuration {
		t.Errorf("SetOptions() TokenDuration = %v, want %v", Options.TokenDuration, options.TokenDuration)
	}
}

func TestAPIKeyAuth_SetStorage(t *testing.T) {
	provider := NewApiKeyAuth().(*APIKeyAuth)
	storage := testutils.NewMockStorage()

	provider.SetStorage(storage)

	if provider.storage == nil {
		t.Errorf("SetStorage() did not set storage")
	}
}

func TestAPIKeyAuth_Decode(t *testing.T) {
	provider := NewApiKeyAuth()

	tests := []struct {
		name        string
		input       string
		expectUser  string
		expectPass  string
		expectError bool
	}{
		{
			name:        "valid api key",
			input:       "test-api-key-12345",
			expectUser:  "exampleUser",     // Based on the validateAPIKey implementation
			expectPass:  "examplePassword", // Based on the validateAPIKey implementation
			expectError: false,
		},
		{
			name:        "any api key works with current implementation",
			input:       "any-key",
			expectUser:  "exampleUser",
			expectPass:  "examplePassword",
			expectError: false,
		},
		{
			name:        "empty api key still works",
			input:       "",
			expectUser:  "exampleUser",
			expectPass:  "examplePassword",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, pass, err := provider.Decode(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("Decode() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Decode() unexpected error = %v", err)
				return
			}

			if !tt.expectError {
				if user != tt.expectUser {
					t.Errorf("Decode() user = %v, want %v", user, tt.expectUser)
				}

				if pass != tt.expectPass {
					t.Errorf("Decode() password = %v, want %v", pass, tt.expectPass)
				}
			}
		})
	}
}

func TestAPIKeyAuth_CreateAccessToken(t *testing.T) {
	provider := NewApiKeyAuth().(*APIKeyAuth)
	storage := testutils.NewMockStorage()
	options := testutils.TestAuthOptions()

	provider.SetStorage(storage)
	provider.SetOptions(options)

	tests := []struct {
		name        string
		userID      string
		apiKey      string
		expectError bool
	}{
		{
			name:        "valid api key",
			userID:      "testuser", // userID is not used in current implementation
			apiKey:      "test-api-key-12345",
			expectError: false,
		},
		{
			name:        "invalid api key",
			userID:      "testuser",
			apiKey:      "invalid-api-key",
			expectError: true,
		},
		{
			name:        "empty api key",
			userID:      "testuser",
			apiKey:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := provider.CreateAccessToken(tt.userID, tt.apiKey)

			if tt.expectError && err == nil {
				t.Errorf("CreateAccessToken() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("CreateAccessToken() unexpected error = %v", err)
				return
			}

			if !tt.expectError {
				if token == nil {
					t.Errorf("CreateAccessToken() returned nil token")
				} else {
					// Test that token has a valid string representation
					tokenString := token.GetToken()
					if tokenString == "" {
						t.Errorf("CreateAccessToken() token string is empty")
					}
				}
			}
		})
	}
}

func TestAPIKeyAuth_CreateRefreshToken(t *testing.T) {
	provider := NewApiKeyAuth().(*APIKeyAuth)
	storage := testutils.NewMockStorage()
	options := testutils.TestAuthOptions()

	provider.SetStorage(storage)
	provider.SetOptions(options)

	tests := []struct {
		name        string
		userID      string
		apiKey      string
		expectError bool
	}{
		{
			name:        "valid api key",
			userID:      "testuser",
			apiKey:      "test-api-key-12345",
			expectError: false,
		},
		{
			name:        "invalid api key",
			userID:      "testuser",
			apiKey:      "invalid-api-key",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := provider.CreateRefreshToken(tt.userID, tt.apiKey)

			if tt.expectError && err == nil {
				t.Errorf("CreateRefreshToken() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("CreateRefreshToken() unexpected error = %v", err)
				return
			}

			if !tt.expectError {
				if token == nil {
					t.Errorf("CreateRefreshToken() returned nil token")
				} else {
					tokenString := token.GetToken()
					if tokenString == "" {
						t.Errorf("CreateRefreshToken() token string is empty")
					}
				}
			}
		})
	}
}

func TestAPIKeyAuth_Validate(t *testing.T) {
	provider := NewApiKeyAuth().(*APIKeyAuth)
	options := testutils.TestAuthOptions()
	provider.SetOptions(options)

	// Test with invalid token
	t.Run("invalid token", func(t *testing.T) {
		_, err := provider.Validate("invalid.jwt.token")
		if err == nil {
			t.Errorf("Validate() expected error with invalid token")
		}
	})

	// Test with empty token
	t.Run("empty token", func(t *testing.T) {
		_, err := provider.Validate("")
		if err == nil {
			t.Errorf("Validate() expected error with empty token")
		}
	})
}

func TestAPIKeyAuth_GrantRefreshToken(t *testing.T) {
	provider := NewApiKeyAuth().(*APIKeyAuth)
	options := testutils.TestAuthOptions()
	provider.SetOptions(options)

	// Test with invalid refresh token
	t.Run("invalid refresh token", func(t *testing.T) {
		_, err := provider.GrantRefreshToken("invalid.refresh.token")
		if err == nil {
			t.Errorf("GrantRefreshToken() expected error with invalid token")
		}
	})
}

func TestValidateAPIKey(t *testing.T) {
	provider := NewApiKeyAuth().(*APIKeyAuth)

	tests := []struct {
		name        string
		input       string
		expectUser  string
		expectPass  string
		expectError bool
	}{
		{
			name:        "valid api key",
			input:       "test-api-key",
			expectUser:  "exampleUser",
			expectPass:  "examplePassword",
			expectError: false,
		},
		{
			name:        "empty api key",
			input:       "",
			expectUser:  "exampleUser",
			expectPass:  "examplePassword",
			expectError: false,
		},
		{
			name:        "long api key",
			input:       "very-long-api-key-with-many-characters",
			expectUser:  "exampleUser",
			expectPass:  "examplePassword",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, pass, err := provider.validateAPIKey(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("validateAPIKey() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("validateAPIKey() unexpected error = %v", err)
				return
			}

			if !tt.expectError {
				if user != tt.expectUser {
					t.Errorf("validateAPIKey() user = %v, want %v", user, tt.expectUser)
				}

				if pass != tt.expectPass {
					t.Errorf("validateAPIKey() password = %v, want %v", pass, tt.expectPass)
				}
			}
		})
	}
}
