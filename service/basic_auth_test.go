package service

import (
	"testing"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/testutils"
)

func TestNewBasicAuth(t *testing.T) {
	provider := NewBasicAuth()

	if provider == nil {
		t.Errorf("NewBasicAuth() returned nil")
	}

	// Test that it implements AuthInterface
	var _ auth.AuthInterface = provider
}

func TestBasicAuth_SetOptions(t *testing.T) {
	provider := NewBasicAuth()
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

func TestBasicAuth_SetStorage(t *testing.T) {
	provider := NewBasicAuth().(*BasicAuth)
	storage := testutils.NewMockStorage()

	provider.SetStorage(storage)

	if provider.storage == nil {
		t.Errorf("SetStorage() did not set storage")
	}
}

func TestBasicAuth_Decode(t *testing.T) {
	provider := NewBasicAuth()

	tests := []struct {
		name        string
		input       string
		expectUser  string
		expectPass  string
		expectError bool
	}{
		{
			name:        "valid basic auth",
			input:       testutils.ValidBasicAuthCredentials(),
			expectUser:  "test@example.com",
			expectPass:  "test-password-hash",
			expectError: false,
		},
		{
			name:        "invalid base64",
			input:       "invalid-base64!@#",
			expectUser:  "",
			expectPass:  "",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectUser:  "",
			expectPass:  "",
			expectError: true,
		},
		{
			name:        "no colon separator",
			input:       "dGVzdA==", // "test" in base64
			expectUser:  "",
			expectPass:  "",
			expectError: true,
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

func TestBasicAuth_CreateAccessToken(t *testing.T) {
	provider := NewBasicAuth().(*BasicAuth)
	storage := testutils.NewMockStorage()
	options := testutils.TestAuthOptions()

	provider.SetStorage(storage)
	provider.SetOptions(options)

	tests := []struct {
		name        string
		userID      string
		hash        string
		expectError bool
	}{
		{
			name:        "valid credentials",
			userID:      "test@example.com",
			hash:        "test-password-hash",
			expectError: false,
		},
		{
			name:        "invalid credentials",
			userID:      "test@example.com",
			hash:        "wrong-password",
			expectError: true,
		},
		{
			name:        "non-existent user",
			userID:      "nonexistent@example.com",
			hash:        "any-password",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := provider.CreateAccessToken(tt.userID, tt.hash)

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

func TestBasicAuth_CreateRefreshToken(t *testing.T) {
	provider := NewBasicAuth().(*BasicAuth)
	storage := testutils.NewMockStorage()
	options := testutils.TestAuthOptions()

	provider.SetStorage(storage)
	provider.SetOptions(options)

	tests := []struct {
		name        string
		userID      string
		hash        string
		expectError bool
	}{
		{
			name:        "valid credentials",
			userID:      "test@example.com",
			hash:        "test-password-hash",
			expectError: false,
		},
		{
			name:        "invalid credentials",
			userID:      "test@example.com",
			hash:        "wrong-password",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := provider.CreateRefreshToken(tt.userID, tt.hash)

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

func TestBasicAuth_Validate(t *testing.T) {
	provider := NewBasicAuth().(*BasicAuth)
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

func TestBasicAuth_GrantRefreshToken(t *testing.T) {
	provider := NewBasicAuth().(*BasicAuth)
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

func TestValidateBasic(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectUser  string
		expectPass  string
		expectError bool
	}{
		{
			name:        "valid credentials",
			input:       "dGVzdEBleGFtcGxlLmNvbTp0ZXN0LXBhc3N3b3JkLWhhc2g=", // test@example.com:test-password-hash
			expectUser:  "test@example.com",
			expectPass:  "test-password-hash",
			expectError: false,
		},
		{
			name:        "invalid base64",
			input:       "invalid-base64!@#",
			expectUser:  "",
			expectPass:  "",
			expectError: true,
		},
		{
			name:        "no colon",
			input:       "dGVzdA==", // "test"
			expectUser:  "",
			expectPass:  "",
			expectError: true,
		},
		{
			name:        "empty after colon",
			input:       "dGVzdDo=", // "test:"
			expectUser:  "test",
			expectPass:  "",
			expectError: false,
		},
		{
			name:        "empty before colon",
			input:       "OnRlc3Q=", // ":test"
			expectUser:  "",
			expectPass:  "test",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, pass, err := validateBasic(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("validateBasic() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("validateBasic() unexpected error = %v", err)
				return
			}

			if !tt.expectError {
				if user != tt.expectUser {
					t.Errorf("validateBasic() user = %v, want %v", user, tt.expectUser)
				}

				if pass != tt.expectPass {
					t.Errorf("validateBasic() password = %v, want %v", pass, tt.expectPass)
				}
			}
		})
	}
}
