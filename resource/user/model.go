package user

import (
	"time"
)

type User struct {
	AccountID uint64
	Name      string
	Mail      string
	Created   uint64
	Access    uint64
	Status    int
	Secret    string
	APIKey    string
	Refresh   string
}

type DTO struct {
	ID        string `json:"id"`
	AccountID uint64 `json:"account_id"`
	Name      string `json:"name"`
	Mail      string `json:"mail"`
	Created   uint64 `json:"created"`
	Access    uint64 `json:"access"`
	Status    int    `json:"status"`
	Secret    string `json:"secret"`
	APIKey    string `json:"apikey"`
	Refresh   string `json:"refresh_token"`
}

type Form struct {
	ID        string `json:"id"`
	AccountID uint64 `json:"account_id"`
	Name      string `json:"name"`
	Mail      string `json:"mail"`
}

func (u *User) ToDto() *DTO {
	return &DTO{
		AccountID: u.AccountID,
		Name:      u.Name,
		Mail:      u.Mail,
		Created:   u.Created,
		Access:    u.Access,
		Status:    u.Status,
		Secret:    u.Secret,
		APIKey:    u.APIKey,
		Refresh:   u.Refresh,
	}
}

func (f *Form) ToModel() *User {
	return &User{
		AccountID: f.AccountID,
		Name:      f.Name,
		Mail:      f.Mail,
		Created:   uint64(time.Now().Unix()),
		Access:    uint64(time.Now().Unix()),
		Status:    1,
	}
}
