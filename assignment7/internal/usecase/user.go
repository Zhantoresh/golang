package usecase

import (
	"fmt"
	"assignment7/internal/entity"
	"assignment7/internal/usecase/repo"
	"assignment7/utils"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, error) {
	return u.repo.RegisterUser(user)
}

func (u *UserUseCase) LoginUser(dto *entity.LoginUserDTO) (string, error) {
	userFromDB, err := u.repo.LoginUser(dto)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(userFromDB.Password, dto.Password) {
		return "", fmt.Errorf("invalid password")
	}

	token, err := utils.GenerateJWT(userFromDB.ID, userFromDB.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.PromoteUser(id)
}