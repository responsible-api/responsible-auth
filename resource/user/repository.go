package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type RepositoryInterface interface {
	Read(mail string) (*User, error)
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Read(username string) (*User, error) {
	user := &User{}
	query := r.db.Table("responsible_api_users").
		Where("mail = ?", username).
		Or("account_id = ?", username)

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
