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

func (r *Repository) Read(username string, hash string) (*User, error) {
	user := &User{}
	// SELECT * FROM responsible_api_users WHERE(mail = ? OR account_id = ?) AND secret = ?
	query := r.db.Table("responsible_api_users").
		Where(r.db.Where("mail = ?", username).Or("account_id = ?", username)).
		Where("secret = ?", hash).
		Limit(1)

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
