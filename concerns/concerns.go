package concerns

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Options struct {
	// Required fields
	SecretKey            string `json:"required"`
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

type Auth struct {
	Options Options
}

type AuthInterface interface {
	GenerateToken(userID string, Hash string) (string, error)
	ValidateToken(token string) (*Claims, error)
}

type Claims struct {
	jwt.RegisteredClaims
	User *User `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
	Audience string `json:"audience,omitempty"`
	Scopes   string `json:"scopes,omitempty"`
}
