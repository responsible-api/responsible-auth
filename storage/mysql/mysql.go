package mysql

import (
	"github.com/responsible-api/responsible-auth/resource/user"
	"github.com/responsible-api/responsible-auth/storage"
	"gorm.io/gorm"
)

// MySQLStorage implements the UserStorage interface using MySQL/GORM
type MySQLStorage struct {
	db *gorm.DB
}

// NewMySQLStorage creates a new MySQL storage implementation
func NewMySQLStorage(db *gorm.DB) storage.UserStorage {
	return &MySQLStorage{
		db: db,
	}
}

// FindUserByCredentials retrieves a user by username/email and validates their credentials
func (m *MySQLStorage) FindUserByCredentials(username, credentials string) (*user.User, error) {
	user := &user.User{}
	query := m.db.Table("responsible_api_users").
		Where(m.db.Where("mail = ?", username).Or("account_id = ?", username)).
		Where("secret = ?", credentials).
		Limit(1)

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByAPIKey retrieves a user by their API key
func (m *MySQLStorage) FindUserByAPIKey(apiKey string) (*user.User, error) {
	user := &user.User{}
	query := m.db.Table("responsible_api_users").
		Where("apikey = ?", apiKey).
		Limit(1)

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateRefreshToken stores a refresh token for a user
func (m *MySQLStorage) UpdateRefreshToken(userID string, refreshToken string) error {
	return m.db.Table("responsible_api_users").
		Where("account_id = ? OR mail = ?", userID, userID).
		Update("refresh_token", refreshToken).Error
}

// ValidateRefreshToken checks if a refresh token is valid for a user
func (m *MySQLStorage) ValidateRefreshToken(refreshToken string) (*user.User, error) {
	user := &user.User{}
	query := m.db.Table("responsible_api_users").
		Where("refresh_token = ?", refreshToken).
		Limit(1)

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
