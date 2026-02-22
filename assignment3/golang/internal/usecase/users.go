package usecase

import (
	"fmt"

	"golang/internal/repository"
	"golang/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) CreateUser(user modules.User) (int, error) {
	id, err := u.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("create user failed: %w", err)
	}
	return id, nil
}

func (u *UserUsecase) UpdateUser(id int, user modules.User) (int64, error) {
	rows, err := u.repo.UpdateUser(id, user)
	if err != nil {
		return 0, fmt.Errorf("update user failed: %w", err)
	}
	return rows, nil
}

func (u *UserUsecase) DeleteUserByID(id int) (int64, error) {
	rows, err := u.repo.DeleteUserByID(id)
	if err != nil {
		return 0, fmt.Errorf("delete user failed: %w", err)
	}
	return rows, nil
}