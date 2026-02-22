package repository

import (
	"github.com/RuslanVNikolov/one-stop-shop/backend/auth-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{ID: id}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Where(&model.User{Email: email}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) SaveRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *UserRepository) FindRefreshToken(tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.Preload("User").
		Where(&model.RefreshToken{TokenHash: tokenHash}).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *UserRepository) RevokeRefreshToken(tokenHash string) error {
	return r.db.Model(&model.RefreshToken{}).
		Where(&model.RefreshToken{TokenHash: tokenHash}).
		Updates(&model.RefreshToken{Revoked: true}).Error
}

func (r *UserRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	return r.db.Model(&model.RefreshToken{}).
		Where(&model.RefreshToken{UserID: userID}).
		Updates(&model.RefreshToken{Revoked: true}).Error
}
