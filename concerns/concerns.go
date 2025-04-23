package concerns

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type ClaimsGeneric struct {
	jwt.RegisteredClaims
	CustomClaims map[string]interface{} `json:"custom,omitempty"`
	Role         string                 `json:"role,omitempty"`
	Scopes       string                 `json:"scopes,omitempty"`
}
