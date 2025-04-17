package internal

import (
	"fmt"
	auth "responsible-api-go"
	"responsible-api-go/concerns"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Grant(userID string, hash string, options auth.AuthOptions) (string, error) {
	if (options.SecretKey == "") || (options.SecretKey == "required") {
		return "", fmt.Errorf("secret key is required")
	}

	// Generate a JWT token by the user ID / username and hash / password
	// Set the expiration time to the specified duration
	// Return the generated token or an error if something goes wrong
	claims := &concerns.ClaimsGeneric{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    setIssuer(options.Issuer),
			IssuedAt:  jwt.NewNumericDate(setIssuedAt(options.IssuedAt)),
			ExpiresAt: jwt.NewNumericDate(setExpiresAt(options.TokenDuration)),
			NotBefore: jwt.NewNumericDate(setNotBefore(options.NotBefore)),
		},
		User: &concerns.ClaimsUser{
			ID: userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(options.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
