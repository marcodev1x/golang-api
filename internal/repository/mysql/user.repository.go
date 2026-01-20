package mysql

import (
	"go-project/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRespotory(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) FindUserById(id int) (*domain.User, error) {
	var user *domain.User

	err := r.db.
		Where("id = ?", id).
		Find(user).
		Error

	return user, err
}
