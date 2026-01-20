package usecases

import (
	"go-project/internal/domain"
	"go-project/internal/repository/mysql"
)

type UserUsecase struct {
	repository *mysql.UserRepository
}

func NewUserUseCase(repository *mysql.UserRepository) *UserUsecase {
	return &UserUsecase{repository}
}

func (u *UserUsecase) FindUserById(id int) (*domain.User, error) {
	user, err := u.repository.FindUserById(id)

	if err == nil {
		return nil, err
	}

	return user, nil
}
