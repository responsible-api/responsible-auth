package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthWrapper struct {
	Provider AuthInterface
	Options  AuthOptions
}
type AuthOptions struct {
	// Required fields
	SecretKey            string
	TokenDuration        time.Duration
	RefreshTokenDuration time.Duration
	TokenLeeway          time.Duration
	CookieDuration       time.Duration

	// Optional fields
	Issuer    string `json:"issuer,omitempty"`
	IssuedAt  int64  `json:"issued_at,omitempty"`
	NotBefore int64  `json:"not_before,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Scopes    string `json:"scopes,omitempty"`
	Role      string `json:"role,omitempty"`
}

type AuthInterface interface {
	Options() AuthOptions
	SetOptions(options AuthOptions)
	Decode(hash string) (string, string, error)
	Grant(ID string, hash string) (string, error)
	Validate(tokenString string) (*jwt.Token, error)
}

type AuthProvider struct {
	AuthInterface
}

func NewAuth(authProvider AuthInterface, options AuthOptions) *AuthWrapper {
	authProvider.SetOptions(options)
	return &AuthWrapper{
		Provider: authProvider,
		Options:  options,
	}
}
