package main

import (
	"log"
	"time"

	auth "responsible-api-go"
	"responsible-api-go/concerns"
)

func main() {
	// c := config.New()
	// fmt.Println("Config loaded:", c)

	// Create a new Auth instance with defined options
	authService := auth.NewAuth(concerns.Options{
		SecretKey:            "your-secret-key",  // Replace with a secure key
		TokenDuration:        5 * time.Minute,    // 5 minute token duration
		RefreshTokenDuration: 24 * time.Hour,     // 1 day refresh token duration
		TokenLeeway:          10 * time.Second,   // 10 seconds leeway
		CookieDuration:       7 * 24 * time.Hour, // 7 days cookie duration
		Issuer:               "your-issuer",      // Replace with your issuer
		IssuedAt:             time.Now().Unix(),  // Time we issued the token, can be now or in the future
		NotBefore:            time.Now().Unix(),  // Time before which the token is not valid
		Subject:              "your-subject",     // Replace with your subject
	})

	// Example user data
	userID := "12345"
	hash := "user-specific-hash"

	// Generate a token for the user
	token, err := authService.GenerateToken(userID, hash)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	_, err = authService.ValidateToken(token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	log.Println("Token is valid.", token)
}
