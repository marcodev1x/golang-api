package rest

import "go-project/internal/usecases"

type UserRest struct {
	usecase *usecases.UserUsecase
}

func NewUserRest(usecases *usecases.UserUsecase) *UserRest {
	return &UserRest{usecases}
}
