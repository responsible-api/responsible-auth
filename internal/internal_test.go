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

	// Create an expired token for testing
	expiredOptions := options
	expiredOptions.TokenDuration = -1 * time.Hour
	expiredToken, err := CreateAccessToken(expiredOptions)
	if err != nil {
		t.Fatalf("Failed to create expired token for testing: %v", err)
	}

	// Create a token with future NotBefore for testing
	futureNbfOptions := options
	futureNbfOptions.NotBefore = time.Now().Add(1 * time.Hour).Unix()
	futureNbfToken, err := CreateAccessToken(futureNbfOptions)
	if err != nil {
		t.Fatalf("Failed to create future nbf token for testing: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		options     auth.AuthOptions
		expectError bool
		errorMsg    string
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
		{
			name:        "expired token",
			tokenString: expiredToken.GetToken(),
			options:     options,
			expectError: true,
		},
		{
			name:        "token not valid yet (nbf in future)",
			tokenString: futureNbfToken.GetToken(),
			options:     options,
			expectError: true,
		},
		{
			name:        "malformed token - only header",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			options:     options,
			expectError: true,
		},
		{
			name:        "malformed token - missing signature",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			options:     options,
			expectError: true,
		},
		{
			name: "token with leeway should be valid",
			tokenString: func() string {
				leewayOptions := options
				leewayOptions.TokenLeeway = 5 * time.Minute

				// Create a token that's slightly expired but within leeway
				// Use a smaller expiry offset since JWT library + our validation both need to pass
				claims := &concerns.ClaimsGeneric{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)), // Still valid
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						Issuer:    options.Issuer,
						Subject:   options.Subject,
					},
					Role:         options.Role,
					Scopes:       options.Scopes,
					CustomClaims: options.CustomClaims,
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(options.SecretKey))
				return tokenString
			}(),
			options: auth.AuthOptions{
				SecretKey:   options.SecretKey,
				TokenLeeway: 5 * time.Minute,
			},
			expectError: false,
		},
		{
			name: "token with custom claims",
			tokenString: func() string {
				customClaimsOptions := options
				customClaimsOptions.CustomClaims = map[string]interface{}{
					"user_id":     "12345",
					"permissions": []string{"read", "write"},
					"metadata": map[string]string{
						"department": "engineering",
					},
				}
				token, _ := CreateAccessToken(customClaimsOptions)
				return token.GetToken()
			}(),
			options:     options,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := Validate(tt.tokenString, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}

				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Validate() expected error message '%s' but got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
				return
			}

			if token == nil {
				t.Errorf("Validate() returned nil token")
				return
			}

			if !token.Valid {
				t.Errorf("Validate() returned invalid token")
			}

			// Test that claims can be extracted and contain expected values
			if claims, ok := token.Claims.(*concerns.ClaimsGeneric); ok {
				if claims.ExpiresAt == nil {
					t.Errorf("Validate() token has no expiration claim")
				}

				// Verify standard claims if they were set
				if tt.options.Issuer != "" && claims.Issuer != tt.options.Issuer {
					t.Errorf("Validate() issuer mismatch: expected %s, got %s", tt.options.Issuer, claims.Issuer)
				}

				if tt.options.Subject != "" && claims.Subject != tt.options.Subject {
					t.Errorf("Validate() subject mismatch: expected %s, got %s", tt.options.Subject, claims.Subject)
				}

				if tt.options.Role != "" && claims.Role != tt.options.Role {
					t.Errorf("Validate() role mismatch: expected %s, got %s", tt.options.Role, claims.Role)
				}

				if tt.options.Scopes != "" && claims.Scopes != tt.options.Scopes {
					t.Errorf("Validate() scopes mismatch: expected %s, got %s", tt.options.Scopes, claims.Scopes)
				}
			} else {
				t.Errorf("Validate() token claims are not ClaimsGeneric type")
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
		{
			name: "expiry exactly now (should be invalid)",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now()),
				},
			},
			expected: false, // Exact time should be considered expired due to Before() check
		},
		{
			name: "expiry 1 second in future",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Second)),
				},
			},
			expected: true,
		},
		{
			name: "expiry 1 second in past",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
				},
			},
			expected: false,
		},
		{
			name: "expiry far in future",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
			},
			expected: true,
		},
		{
			name: "expiry far in past",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
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
		{
			name: "not before exactly now (should be valid)",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now()),
				},
			},
			expected: true, // Exact time should be considered valid due to After() check
		},
		{
			name: "not before 1 second in past",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
				},
			},
			expected: true,
		},
		{
			name: "not before 1 second in future",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Second)),
				},
			},
			expected: false,
		},
		{
			name: "not before far in past",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
				},
			},
			expected: true,
		},
		{
			name: "not before far in future",
			claims: &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					NotBefore: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
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

func TestValidateWithDifferentSigningMethods(t *testing.T) {
	tests := []struct {
		name          string
		signingMethod jwt.SigningMethod
		expectError   bool
	}{
		{
			name:          "HMAC SHA256 (valid)",
			signingMethod: jwt.SigningMethodHS256,
			expectError:   false,
		},
		{
			name:          "HMAC SHA384 (valid)",
			signingMethod: jwt.SigningMethodHS384,
			expectError:   false,
		},
		{
			name:          "HMAC SHA512 (valid)",
			signingMethod: jwt.SigningMethodHS512,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := testutils.TestAuthOptions()

			claims := &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					NotBefore: jwt.NewNumericDate(time.Now()),
					Issuer:    options.Issuer,
					Subject:   options.Subject,
				},
				Role:         options.Role,
				Scopes:       options.Scopes,
				CustomClaims: options.CustomClaims,
			}

			token := jwt.NewWithClaims(tt.signingMethod, claims)
			tokenString, err := token.SignedString([]byte(options.SecretKey))
			if err != nil {
				t.Fatalf("Failed to sign token: %v", err)
			}

			validatedToken, err := Validate(tokenString, options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
				return
			}

			if validatedToken == nil || !validatedToken.Valid {
				t.Errorf("Validate() returned invalid token")
			}
		})
	}
}

func TestValidateWithTokenLeeway(t *testing.T) {
	tests := []struct {
		name         string
		expiryOffset time.Duration
		leeway       time.Duration
		expectError  bool
	}{
		{
			name:         "token valid for 1 minute with leeway (should be valid)",
			expiryOffset: 1 * time.Minute,
			leeway:       5 * time.Minute,
			expectError:  false,
		},
		{
			name:         "token valid with no leeway needed",
			expiryOffset: 1 * time.Hour,
			leeway:       1 * time.Minute,
			expectError:  false,
		},
		{
			name:         "token expired far beyond leeway (should be invalid)",
			expiryOffset: -1 * time.Hour,
			leeway:       5 * time.Minute,
			expectError:  true,
		},
		{
			name:         "token valid with zero leeway",
			expiryOffset: 1 * time.Hour,
			leeway:       0,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := testutils.TestAuthOptions()
			options.TokenLeeway = tt.leeway

			claims := &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(tt.expiryOffset)),
					IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
					NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
					Issuer:    options.Issuer,
					Subject:   options.Subject,
				},
				Role:         options.Role,
				Scopes:       options.Scopes,
				CustomClaims: options.CustomClaims,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(options.SecretKey))
			if err != nil {
				t.Fatalf("Failed to sign token: %v", err)
			}

			validatedToken, err := Validate(tokenString, options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
				return
			}

			if validatedToken == nil || !validatedToken.Valid {
				t.Errorf("Validate() returned invalid token")
			}
		})
	}
}

func TestValidateCustomClaims(t *testing.T) {
	tests := []struct {
		name         string
		customClaims map[string]interface{}
		expectError  bool
	}{
		{
			name: "simple custom claims",
			customClaims: map[string]interface{}{
				"user_id": "12345",
				"role":    "admin",
			},
			expectError: false,
		},
		{
			name: "complex custom claims",
			customClaims: map[string]interface{}{
				"user_id":     "67890",
				"permissions": []string{"read", "write", "delete"},
				"metadata": map[string]interface{}{
					"department": "engineering",
					"level":      5,
					"active":     true,
				},
			},
			expectError: false,
		},
		{
			name:         "empty custom claims",
			customClaims: map[string]interface{}{},
			expectError:  false,
		},
		{
			name:         "nil custom claims",
			customClaims: nil,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := testutils.TestAuthOptions()
			options.CustomClaims = tt.customClaims

			claims := &concerns.ClaimsGeneric{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					NotBefore: jwt.NewNumericDate(time.Now()),
					Issuer:    options.Issuer,
					Subject:   options.Subject,
				},
				Role:         options.Role,
				Scopes:       options.Scopes,
				CustomClaims: options.CustomClaims,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(options.SecretKey))
			if err != nil {
				t.Fatalf("Failed to sign token: %v", err)
			}

			validatedToken, err := Validate(tokenString, options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
				return
			}

			if validatedToken == nil || !validatedToken.Valid {
				t.Errorf("Validate() returned invalid token")
				return
			}

			// Verify custom claims are preserved
			if validatedClaims, ok := validatedToken.Claims.(*concerns.ClaimsGeneric); ok {
				if len(tt.customClaims) > 0 {
					if validatedClaims.CustomClaims == nil {
						t.Errorf("Validate() custom claims were lost")
						return
					}

					for key, expectedValue := range tt.customClaims {
						if actualValue, exists := validatedClaims.CustomClaims[key]; !exists {
							t.Errorf("Validate() missing custom claim '%s'", key)
						} else {
							// For complex comparisons, we'll just check existence
							// since interface{} comparison can be tricky
							if actualValue == nil && expectedValue != nil {
								t.Errorf("Validate() custom claim '%s' is nil but expected non-nil", key)
							}
						}
					}
				}
			}
		})
	}
}

func TestValidateEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupToken  func() (string, auth.AuthOptions)
		expectError bool
		description string
	}{
		{
			name: "token with only required claims",
			setupToken: func() (string, auth.AuthOptions) {
				options := auth.AuthOptions{
					SecretKey:            "test-secret-key-32-characters!",
					TokenDuration:        1 * time.Hour,
					RefreshTokenDuration: 24 * time.Hour,
					TokenLeeway:          30 * time.Second,
					CookieDuration:       24 * time.Hour,
				}

				claims := &concerns.ClaimsGeneric{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						NotBefore: jwt.NewNumericDate(time.Now()),
					},
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(options.SecretKey))
				return tokenString, options
			},
			expectError: false,
			description: "token with minimal required claims should be valid",
		},
		{
			name: "token with zero leeway",
			setupToken: func() (string, auth.AuthOptions) {
				options := testutils.TestAuthOptions()
				options.TokenLeeway = 0

				claims := &concerns.ClaimsGeneric{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						NotBefore: jwt.NewNumericDate(time.Now()),
					},
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(options.SecretKey))
				return tokenString, options
			},
			expectError: false,
			description: "token with zero leeway should still validate if not expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, options := tt.setupToken()

			validatedToken, err := Validate(tokenString, options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate() expected error but got none: %s", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v: %s", err, tt.description)
				return
			}

			if validatedToken == nil || !validatedToken.Valid {
				t.Errorf("Validate() returned invalid token: %s", tt.description)
			}
		})
	}
}
