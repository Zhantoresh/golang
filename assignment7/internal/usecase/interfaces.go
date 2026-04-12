package usecase

import "assignment7/internal/entity"

type UserInterface interface {
	RegisterUser(user *entity.User) (*entity.User, error)
	LoginUser(dto *entity.LoginUserDTO) (string, error)
	GetUserByID(id string) (*entity.User, error)
	PromoteUser(id string) error
}