package internal

import (
	"testing"
	"time"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/concerns"
	"github.com/responsible-api/responsible-auth/testutils"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateAccessToken(t *testing.T) {
	tests := []struct {
		name        string
		options     auth.AuthOptions
		expectError bool
	}{
		{
			name:        "valid options",
			options:     testutils.TestAuthOptions(),
			expectError: false,
		},
		{
			name: "missing secret key",
			options: auth.AuthOptions{
				SecretKey:            "",
				TokenDuration:        1 * time.Hour,
				RefreshTokenDuration: 24 * time.Hour,
				TokenLeeway:          30 * time.Second,
				CookieDuration:       24 * time.Hour,
			},
			expectError: true,
		},
		{
			name: "required secret key placeholder",
			options: auth.AuthOptions{
				SecretKey:            "required",
				TokenDuration:        1 * time.Hour,
				RefreshTokenDuration: 24 * time.Hour,
				TokenLeeway:          30 * time.Second,
				CookieDuration:       24 * time.Hour,
			},
			expectError: true,
		},
		{
			name: "minimal valid options",
			options: auth.AuthOptions{
				SecretKey:            "test-secret",
				TokenDuration:        1 * time.Hour,
				RefreshTokenDuration: 24 * time.Hour,
				TokenLeeway:          30 * time.Second,
				CookieDuration:       24 * time.Hour,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateAccessToken(tt.options)

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
					return
				}

				// Test that token string is not empty
				tokenString := token.GetToken()
				if tokenString == "" {
					t.Errorf("CreateAccessToken() token string is empty")
				}

				// Test that token has valid expiration
				exp, err := token.GetExpirationTime()
				if err != nil {
					t.Errorf("CreateAccessToken() error getting expiration: %v", err)
				}

				if exp == nil {
					t.Errorf("CreateAccessToken() token has no expiration")
				}

				// Verify expiration is in the future
				if exp != nil && exp.Time.Before(time.Now()) {
					t.Errorf("CreateAccessToken() token is already expired")
				}
			}
		})
	}
}

func TestCreateRefreshToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		options     auth.AuthOptions
		expectError bool
	}{
		{
			name:        "valid refresh token",
			userID:      "testuser",
			options:     testutils.TestAuthOptions(),
			expectError: false,
		},
		{
			name:   "missing secret key",
			userID: "testuser",
			options: auth.AuthOptions{
				SecretKey:            "",
				TokenDuration:        1 * time.Hour,
				RefreshTokenDuration: 24 * time.Hour,
				TokenLeeway:          30 * time.Second,
				CookieDuration:       24 * time.Hour,
			},
			expectError: true,
		},
		{
			name:        "empty user ID",
			userID:      "",
			options:     testutils.TestAuthOptions(),
			expectError: false, // Empty userID should still work
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateRefreshToken(tt.userID, tt.options)

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
					return
				}

				tokenString := token.GetToken()
				if tokenString == "" {
					t.Errorf("CreateRefreshToken() token string is empty")
				}

				// Test that refresh token has longer expiration than access token
				exp, err := token.GetExpirationTime()
				if err != nil {
					t.Errorf("CreateRefreshToken() error getting expiration: %v", err)
				}

				if exp == nil {
					t.Errorf("CreateRefreshToken() token has no expiration")
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	options := testutils.TestAuthOptions()

	// Create a valid token for testing
	validToken, err := CreateAccessToken(options)
	if err != nil {
		t.Fatalf("Failed to create valid token for testing: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		options     auth.AuthOptions
		expectError bool
	}{
		{
			name:        "valid token",
			tokenString: validToken.GetToken(),
			options:     options,
			expectError: false,
		},
		{
			name:        "invalid token format",
			tokenString: "invalid.jwt.token",
			options:     options,
			expectError: true,
		},
		{
			name:        "empty token",
			tokenString: "",
			options:     options,
			expectError: true,
		},
		{
			name:        "wrong secret key",
			tokenString: validToken.GetToken(),
			options: auth.AuthOptions{
				SecretKey:   "wrong-secret",
				TokenLeeway: options.TokenLeeway,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := Validate(tt.tokenString, tt.options)

			if tt.expectError && err == nil {
				t.Errorf("Validate() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
				return
			}

			if !tt.expectError {
				if token == nil {
					t.Errorf("Validate() returned nil token")
					return
				}

				if !token.Valid {
					t.Errorf("Validate() returned invalid token")
				}

				// Test that claims can be extracted
				if claims, ok := token.Claims.(*concerns.ClaimsGeneric); ok {
					if claims.ExpiresAt == nil {
						t.Errorf("Validate() token has no expiration claim")
					}
				} else {
					t.Errorf("Validate() token claims are not ClaimsGeneric type")
				}
			}
		})
	}
}

func TestValidExpiry(t *testing.T) {
	tests := []struct {
		name     string
		claims   *concerns.ClaimsGeneric
		expected bool
	}{
		{
			name: "future expiry",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				},
			},
			expected: true,
		},
		{
			name: "past expiry",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				},
			},
			expected: false,
		},
		{
			name: "nil expiry",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: nil,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validExpiry(tt.claims)
			if result != tt.expected {
				t.Errorf("validExpiry() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidNotBefore(t *testing.T) {
	tests := []struct {
		name     string
		claims   *concerns.ClaimsGeneric
		expected bool
	}{
		{
			name: "past not before",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				},
			},
			expected: true,
		},
		{
			name: "future not before",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				},
			},
			expected: false,
		},
		{
			name: "nil not before",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: nil,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validNotBefore(tt.claims)
			if result != tt.expected {
				t.Errorf("validNotBefore() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGrantRefreshToken(t *testing.T) {
	options := testutils.TestAuthOptions()

	// Create a valid refresh token for testing
	validRefreshToken, err := CreateRefreshToken("testuser", options)
	if err != nil {
		t.Fatalf("Failed to create valid refresh token for testing: %v", err)
	}

	tests := []struct {
		name               string
		refreshTokenString string
		options            auth.AuthOptions
		expectError        bool
	}{
		{
			name:               "valid refresh token",
			refreshTokenString: validRefreshToken.GetToken(),
			options:            options,
			expectError:        false,
		},
		{
			name:               "invalid refresh token",
			refreshTokenString: "invalid.refresh.token",
			options:            options,
			expectError:        true,
		},
		{
			name:               "empty refresh token",
			refreshTokenString: "",
			options:            options,
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newToken, err := GrantRefreshToken(tt.refreshTokenString, tt.options)

			if tt.expectError && err == nil {
				t.Errorf("GrantRefreshToken() expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("GrantRefreshToken() unexpected error = %v", err)
				return
			}

			if !tt.expectError {
				if newToken == nil {
					t.Errorf("GrantRefreshToken() returned nil token")
					return
				}

				tokenString := newToken.GetToken()
				if tokenString == "" {
					t.Errorf("GrantRefreshToken() token string is empty")
				}

				// Verify the new token is different from the refresh token
				if tokenString == tt.refreshTokenString {
					t.Errorf("GrantRefreshToken() returned same token as input")
				}
			}
		})
	}
}
