package main

import (
	"log"
	"time"

	"github.com/responsible-api/responsible-auth/auth"
	"github.com/responsible-api/responsible-auth/examples/memory"
	"github.com/responsible-api/responsible-auth/resource/access"
	"github.com/responsible-api/responsible-auth/service"
)

// Access token example
func main() {
	// Create in-memory storage implementation
	storage := memory.NewInMemoryStorage()

	// Create auth service with in-memory storage
	authService := auth.NewAuth(service.NewBasicAuth(), storage, auth.AuthOptions{
		SecretKey:            "8m$~t^GbEW<<>cE$BWr5m>)rA>ifVa(3", // Replace with a secure key
		TokenDuration:        5 * time.Hour,                      // 5 minute token duration
		RefreshTokenDuration: 24 * 7 * time.Hour,                 // 7 day refresh token duration
		TokenLeeway:          10 * time.Second,                   // 10 seconds leeway
		CookieDuration:       7 * 24 * time.Hour,                 // 7 days cookie duration
		Issuer:               "https://example.com",              // Replace with your issuer
		IssuedAt:             time.Now().Unix(),                  // Time we issued the token, can be now or in the future
		NotBefore:            time.Now().Unix(),                  // Time before which the token is not valid
		Subject:              "test-user",                        // Replace with your subject,

		// Custom claims if needed
		CustomClaims: map[string]interface{}{
			"organization": "example-org",
			"tier":         "premium",
		},
	})

	// Decode the basic auth credentials (test@example.com:ipHEh|$==*#59@|ftT;IER^qgGG_sz!w)
	user, pass, err := authService.Provider.Decode("dGVzdEBleGFtcGxlLmNvbTppcEhFaHwkPT0qIzU5QHxmdFQ7SUVSXnFnR0dfc3oidw==")
	if err != nil {
		log.Fatalf("Failed to decode basic auth: %v", err)
	}
	log.Printf("Decoded user: %s", user)

	// Grant a token for the user
	token, err := authService.Provider.CreateAccessToken(user, pass)
	if err != nil {
		log.Fatalf("Failed to grant token: %v", err)
	}

	// Create a new model instance and set the values
	expiry, err := token.GetExpirationTime()
	if err != nil {
		log.Println("Failed to get expiration time", err)
	}

	refreshToken, err := authService.Provider.CreateRefreshToken(user, pass)
	if err != nil {
		log.Fatalf("Failed to create refresh token: %v", err)
	}

	model := access.NewModel()
	model.WithAccessToken(token.GetToken())
	model.WithRefreshToken(refreshToken.GetToken())
	model.WithExpiresIn(expiry.Unix())
	model.WithCreatedAt(time.Now().Unix())

	response := model.ToResponseDTO()
	log.Printf("Access Token: %s", response.AccessToken)
	log.Printf("Refresh Token: %s", response.RefreshToken)
	log.Println("🎉 Successfully generated tokens using in-memory storage!")
}
