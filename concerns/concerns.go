package concerns

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type ClaimsGeneric struct {
	RegisteredClaims
	CustomClaims map[string]interface{} `json:"custom,omitempty"`
	Role         string                 `json:"role,omitempty"`
	Scopes       string                 `json:"scopes,omitempty"`
}
