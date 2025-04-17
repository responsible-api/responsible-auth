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
	User *ClaimsUser `json:"user"`
}

type ClaimsUser struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
	Audience string `json:"audience,omitempty"`
	Scopes   string `json:"scopes,omitempty"`
}
