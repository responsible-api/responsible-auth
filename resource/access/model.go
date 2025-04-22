package access

import (
	"strings"
	"time"
)

type Model struct {
	DTO *ResponseDTO
}

type ResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scopes       string `json:"scope,omitempty"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at,omitempty"`
}

func NewModel() *Model {
	return &Model{
		DTO: &ResponseDTO{},
	}
}

func (b *Model) WithAccessToken(token string) {
	b.DTO.AccessToken = token
}

func (b *Model) WithRefreshToken(token string) {
	b.DTO.RefreshToken = token
}

func (b *Model) WithExpiresIn(expiresIn int64) {
	b.DTO.ExpiresIn = expiresIn - time.Now().Unix()
}

func (b *Model) WithCreatedAt(createdAt int64) {
	b.DTO.CreatedAt = createdAt
}

func (b *Model) WithUpdatedAt(updatedAt int64) {
	b.DTO.UpdatedAt = updatedAt
}

func (b *Model) WithScopesString(scopes string) {
	if strings.TrimSpace(scopes) == "" {
		return
	}
	b.DTO.Scopes = strings.TrimSpace(scopes)
}

func (b *Model) WithScopes(scopes []string) {
	if len(scopes) == 0 {
		return
	}
	b.DTO.Scopes = strings.Join(scopes, " ")
}

// Transform the DTO set to a response DTO
func (b *Model) ToResponseDTO() *ResponseDTO {
	return &ResponseDTO{
		AccessToken:  b.DTO.AccessToken,
		RefreshToken: b.DTO.RefreshToken,
		ExpiresIn:    b.DTO.ExpiresIn,
		Scopes:       b.DTO.Scopes,
		CreatedAt:    b.DTO.CreatedAt,
		UpdatedAt:    b.DTO.UpdatedAt,
	}
}
