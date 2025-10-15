package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/responsible-api/responsible-auth/resource/access"
	"github.com/responsible-api/responsible-auth/storage"
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

	// Custom claims
	CustomClaims map[string]interface{} `json:"custom_claims,omitempty"`
}

type AuthInterface interface {
	Options() AuthOptions
	SetOptions(options AuthOptions)
	SetStorage(storage storage.UserStorage)
	Decode(hash string) (string, string, error)
	CreateAccessToken(userID string, hash string) (*access.RToken, error)
	CreateRefreshToken(userID string, hash string) (*access.RToken, error)
	GrantRefreshToken(refreshTokenString string) (*access.RToken, error)
	Validate(tokenString string) (*jwt.Token, error)
}

type AuthProvider struct {
	AuthInterface
}

func NewAuth(authProvider AuthInterface, storage storage.UserStorage, options AuthOptions) *AuthWrapper {
	authProvider.SetOptions(options)
	authProvider.SetStorage(storage)
	return &AuthWrapper{
		Provider: authProvider,
		Options:  options,
	}
}
