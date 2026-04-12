package repo

import (
	"fmt"
	"assignment7/internal/entity"
	"assignment7/pkg/postgres"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	if err := r.PG.Conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) LoginUser(dto *entity.LoginUserDTO) (*entity.User, error) {
	var user entity.User
	if err := r.PG.Conn.Where("username = ?", dto.Username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *UserRepo) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	if err := r.PG.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) PromoteUser(id string) error {
	return r.PG.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", "admin").Error
}