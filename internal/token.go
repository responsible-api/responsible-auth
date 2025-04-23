package internal

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vince-scarpa/responsible-api-go/auth"
	"github.com/vince-scarpa/responsible-api-go/concerns"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(options auth.AuthOptions) (string, error) {
	if (options.SecretKey == "") || (options.SecretKey == "required") {
		return "", fmt.Errorf("secret key is required")
	}

	// Generate a JWT token via the supplied options set
	// Set the expiration time to the specified duration
	// Return the generated token or an error if something goes wrong
	claims := &concerns.ClaimsGeneric{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    setIssuer(options.Issuer),
			Subject:   setSubject(options.Subject),
			IssuedAt:  jwt.NewNumericDate(setIssuedAt(options.IssuedAt)),
			ExpiresAt: jwt.NewNumericDate(setExpiresAt(options.TokenDuration)),
			NotBefore: jwt.NewNumericDate(setNotBefore(options.NotBefore)),
		},
		Role:   options.Role,
		Scopes: options.Scopes,
		// Custom claims can be added here
		CustomClaims: options.CustomClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(options.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(username string, options auth.AuthOptions) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(options.RefreshTokenDuration).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(options.SecretKey))
	if err != nil {
		return "", err
	}

	// Return the refresh token string
	return refreshTokenString, nil
}

func GrantRefreshToken(refreshTokenString string, options auth.AuthOptions) (string, error) {
	// Parse and verify the requested refresh token to grant a new access token
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(options.SecretKey), nil
	})

	if err != nil || !refreshToken.Valid {
		log.Println("Error parsing refresh token:", err)
		return "", fmt.Errorf("invalid refresh token")
	}

	// Generate a new access token if refresh token is valid
	if _, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		newAccessToken, err := CreateAccessToken(options)
		if err != nil {
			return "", err
		}
		return newAccessToken, nil
	}

	// If the refresh token is not valid, return an error
	return "", fmt.Errorf("invalid refresh token")
}

// setIssuer sets the issuer for the token.
// If issuer in the option is empty or doesn't exist, then default to "default-issuer".
// Otherwise, it sets the issuer to the requested value.
// Options.Issuer string `json:"issuer,omitempty"`
func setIssuer(issuer string) string {
	if issuer == "" {
		return "default-issuer"
	}
	return issuer
}

// setIssuedAt sets the issued at time for the token.
// If issuedAt in the option is zero or doesn't exist, then default to current time.
// Otherwise, it sets the issued at time to the requested time.
// Options.IssuedAt int64 `json:"issued_at,omitempty"`
func setIssuedAt(issuedAt int64) time.Time {
	if issuedAt == 0 {
		// Default to current time if not set
		return time.Now()
	}
	return time.Unix(int64(issuedAt), 0)
}

// setExpiresAt sets the expiration time for the token.
// If expiresAtDuration in the option is zero or doesn't exist, then default to 15 minutes.
// Otherwise, it sets the expiration time to the current time plus the requested duration.
// Options.TokenDuration time.Duration
func setExpiresAt(expiresAtDuration time.Duration) time.Time {
	if expiresAtDuration == 0 {
		// Default to 15 minutes if not set
		return time.Now().Add(15 * time.Minute)
	}
	return time.Now().Add(expiresAtDuration)
}

// setNotBefore sets the not before time for the token.
// If notBefore in the option is zero or doesn't exist, then default to current time.
// Otherwise, it sets the not before time to the requested time.
// Options.NotBefore int64 `json:"not_before,omitempty"`
func setNotBefore(notBefore int64) time.Time {
	if notBefore == 0 {
		// Default to current time if not set
		return time.Now()
	}
	return time.Unix(int64(notBefore), 0)
}

// setSubject sets the subject for the token.
// If subject in the option is empty or doesn't exist, then default to nil.
// Otherwise, it sets the subject to the requested value.
// Options.Subject string `json:"subject,omitempty"`
func setSubject(subject string) string {
	if subject == "" {
		return ""
	}
	return subject
}
