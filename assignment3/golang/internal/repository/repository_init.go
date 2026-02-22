package repository

import (
	"golang/internal/repository/_postgres"
	"golang/internal/repository/_postgres/users"
)

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}