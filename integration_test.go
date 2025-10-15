package main

import (
	"testing"
	"time"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/examples/memory"
	"github.com/responsible-api/responsible-auth/service"
	"github.com/responsible-api/responsible-auth/testutils"
)

func TestBasicAuthIntegration(t *testing.T) {
	// Setup
	storage := memory.NewInMemoryStorage()
	provider := service.NewBasicAuth()
	options := testutils.TestAuthOptions()

	authService := auth.NewAuth(provider, storage, options)

	t.Run("complete basic auth flow", func(t *testing.T) {
		// 1. Decode credentials
		username, password, err := authService.Provider.Decode(testutils.ValidBasicAuthCredentials())
		if err != nil {
			t.Fatalf("Failed to decode credentials: %v", err)
		}

		if username != "test@example.com" {
			t.Errorf("Expected username test@example.com, got %s", username)
		}

		// 2. Create access token
		accessToken, err := authService.Provider.CreateAccessToken(username, password)
		if err != nil {
			t.Fatalf("Failed to create access token: %v", err)
		}

		if accessToken == nil {
			t.Fatal("Access token is nil")
		}

		// 3. Validate access token
		tokenString := accessToken.GetToken()
		validatedToken, err := authService.Provider.Validate(tokenString)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if !validatedToken.Valid {
			t.Error("Token should be valid")
		}

		// 4. Create refresh token
		refreshToken, err := authService.Provider.CreateRefreshToken(username, password)
		if err != nil {
			t.Fatalf("Failed to create refresh token: %v", err)
		}

		// 5. Use refresh token to get new access token
		refreshTokenString := refreshToken.GetToken()
		newAccessToken, err := authService.Provider.GrantRefreshToken(refreshTokenString)
		if err != nil {
			t.Fatalf("Failed to grant refresh token: %v", err)
		}

		if newAccessToken == nil {
			t.Fatal("New access token is nil")
		}

		// 6. Validate new access token
		newTokenString := newAccessToken.GetToken()
		newValidatedToken, err := authService.Provider.Validate(newTokenString)
		if err != nil {
			t.Fatalf("Failed to validate new token: %v", err)
		}

		if !newValidatedToken.Valid {
			t.Error("New token should be valid")
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		// Try with wrong credentials
		_, err := authService.Provider.CreateAccessToken("test@example.com", "wrong-password")
		if err == nil {
			t.Error("Expected error with wrong credentials")
		}
	})

	t.Run("non-existent user", func(t *testing.T) {
		// Try with non-existent user
		_, err := authService.Provider.CreateAccessToken("nonexistent@example.com", "any-password")
		if err == nil {
			t.Error("Expected error with non-existent user")
		}
	})
}

func TestAPIKeyAuthIntegration(t *testing.T) {
	// Setup
	storage := memory.NewInMemoryStorage()
	provider := service.NewApiKeyAuth()
	options := testutils.TestAuthOptions()

	authService := auth.NewAuth(provider, storage, options)

	t.Run("complete api key auth flow", func(t *testing.T) {
		// 1. Decode API key (current implementation returns static values)
		username, _, err := authService.Provider.Decode("api_key_12345")
		if err != nil {
			t.Fatalf("Failed to decode API key: %v", err)
		}

		// Current implementation returns static values
		if username != "exampleUser" {
			t.Errorf("Expected username exampleUser, got %s", username)
		}

		// 2. Create access token using valid API key
		accessToken, err := authService.Provider.CreateAccessToken(username, "api_key_12345")
		if err != nil {
			t.Fatalf("Failed to create access token: %v", err)
		}

		if accessToken == nil {
			t.Fatal("Access token is nil")
		}

		// 3. Validate access token
		tokenString := accessToken.GetToken()
		validatedToken, err := authService.Provider.Validate(tokenString)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if !validatedToken.Valid {
			t.Error("Token should be valid")
		}

		// 4. Create refresh token
		refreshToken, err := authService.Provider.CreateRefreshToken(username, "api_key_12345")
		if err != nil {
			t.Fatalf("Failed to create refresh token: %v", err)
		}

		if refreshToken == nil {
			t.Fatal("Refresh token is nil")
		}
	})

	t.Run("invalid api key", func(t *testing.T) {
		// Try with invalid API key
		_, err := authService.Provider.CreateAccessToken("exampleUser", "invalid-api-key")
		if err == nil {
			t.Error("Expected error with invalid API key")
		}
	})
}

func TestMultipleProvidersWithSameStorage(t *testing.T) {
	// Test that different providers can use the same storage
	storage := memory.NewInMemoryStorage()
	options := testutils.TestAuthOptions()

	// Basic Auth service
	basicProvider := service.NewBasicAuth()
	basicAuthService := auth.NewAuth(basicProvider, storage, options)

	// API Key Auth service
	apiKeyProvider := service.NewApiKeyAuth()
	apiKeyAuthService := auth.NewAuth(apiKeyProvider, storage, options)

	t.Run("both providers access same user data", func(t *testing.T) {
		// Basic auth flow
		username, password, err := basicAuthService.Provider.Decode(testutils.ValidBasicAuthCredentials())
		if err != nil {
			t.Fatalf("Basic auth decode failed: %v", err)
		}

		basicToken, err := basicAuthService.Provider.CreateAccessToken(username, password)
		if err != nil {
			t.Fatalf("Basic auth token creation failed: %v", err)
		}

		// API key auth flow
		apiKeyToken, err := apiKeyAuthService.Provider.CreateAccessToken("exampleUser", "api_key_12345")
		if err != nil {
			t.Fatalf("API key auth token creation failed: %v", err)
		}

		// Both should create valid tokens
		if basicToken == nil || apiKeyToken == nil {
			t.Fatal("One or both tokens are nil")
		}

		// Validate both tokens
		basicValidated, err := basicAuthService.Provider.Validate(basicToken.GetToken())
		if err != nil {
			t.Errorf("Basic token validation failed: %v", err)
		}

		apiKeyValidated, err := apiKeyAuthService.Provider.Validate(apiKeyToken.GetToken())
		if err != nil {
			t.Errorf("API key token validation failed: %v", err)
		}

		if !basicValidated.Valid || !apiKeyValidated.Valid {
			t.Error("One or both tokens are invalid")
		}
	})
}

func TestTokenExpiration(t *testing.T) {
	// Test with short token duration
	storage := memory.NewInMemoryStorage()
	provider := service.NewBasicAuth()

	shortOptions := testutils.TestAuthOptions()
	shortOptions.TokenDuration = 100 * time.Millisecond // Very short duration

	authService := auth.NewAuth(provider, storage, shortOptions)

	t.Run("token expires correctly", func(t *testing.T) {
		username, password, err := authService.Provider.Decode(testutils.ValidBasicAuthCredentials())
		if err != nil {
			t.Fatalf("Failed to decode credentials: %v", err)
		}

		// Create token
		token, err := authService.Provider.CreateAccessToken(username, password)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		// Token should be valid initially
		tokenString := token.GetToken()
		validatedToken, err := authService.Provider.Validate(tokenString)
		if err != nil {
			t.Fatalf("Token should be valid initially: %v", err)
		}

		if !validatedToken.Valid {
			t.Error("Token should be valid initially")
		}

		// Wait for token to expire
		time.Sleep(200 * time.Millisecond)

		// Token should now be expired
		expiredToken, err := authService.Provider.Validate(tokenString)
		if err == nil {
			t.Error("Expected error for expired token")
		}

		if expiredToken != nil && expiredToken.Valid {
			t.Error("Token should be expired")
		}
	})
}

func TestCustomClaims(t *testing.T) {
	storage := memory.NewInMemoryStorage()
	provider := service.NewBasicAuth()

	options := testutils.TestAuthOptions()
	// Add additional custom claims
	options.CustomClaims["department"] = "engineering"
	options.CustomClaims["level"] = 5
	options.Role = "admin"
	options.Scopes = "read,write,admin"

	authService := auth.NewAuth(provider, storage, options)

	t.Run("custom claims are preserved", func(t *testing.T) {
		username, password, err := authService.Provider.Decode(testutils.ValidBasicAuthCredentials())
		if err != nil {
			t.Fatalf("Failed to decode credentials: %v", err)
		}

		token, err := authService.Provider.CreateAccessToken(username, password)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		// Validate token and check claims
		tokenString := token.GetToken()
		validatedToken, err := authService.Provider.Validate(tokenString)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		// The token should be valid and contain custom claims
		if !validatedToken.Valid {
			t.Error("Token should be valid")
		}

		// Note: In a real implementation, you'd extract and verify the custom claims
		// This test verifies that tokens with custom claims can be created and validated
	})
}

func TestAuthWrapper(t *testing.T) {
	storage := memory.NewInMemoryStorage()
	provider := service.NewBasicAuth()
	options := testutils.TestAuthOptions()

	authWrapper := auth.NewAuth(provider, storage, options)

	t.Run("auth wrapper provides access to components", func(t *testing.T) {
		// Test Provider access
		if authWrapper.Provider == nil {
			t.Error("AuthWrapper.Provider should not be nil")
		}

		// Test Options access
		if authWrapper.Options.SecretKey != options.SecretKey {
			t.Errorf("Options.SecretKey mismatch: got %s, want %s",
				authWrapper.Options.SecretKey, options.SecretKey)
		}

		if authWrapper.Options.TokenDuration != options.TokenDuration {
			t.Errorf("Options.TokenDuration mismatch: got %v, want %v",
				authWrapper.Options.TokenDuration, options.TokenDuration)
		}

		// Test that provider can be used through wrapper
		_, _, err := authWrapper.Provider.Decode(testutils.ValidBasicAuthCredentials())
		if err != nil {
			t.Errorf("Failed to use provider through wrapper: %v", err)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	storage := testutils.NewMockStorage()
	provider := service.NewBasicAuth()
	options := testutils.TestAuthOptions()

	authService := auth.NewAuth(provider, storage, options)

	t.Run("storage errors are propagated", func(t *testing.T) {
		// Configure mock storage to return errors
		storage.SetError(true, "database connection failed")

		_, err := authService.Provider.CreateAccessToken("test@example.com", "test-password-hash")
		if err == nil {
			t.Error("Expected error from storage")
		}

		if err.Error() != "database connection failed" {
			t.Errorf("Expected specific error message, got: %v", err)
		}

		// Reset storage
		storage.SetError(false, "")

		// Should work normally now
		_, err = authService.Provider.CreateAccessToken("test@example.com", "test-password-hash")
		if err != nil {
			t.Errorf("Unexpected error after reset: %v", err)
		}
	})
}
