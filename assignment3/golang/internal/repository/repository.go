package repository

import "golang/pkg/modules"

type UserRepository interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(u modules.User) (int, error)
	UpdateUser(id int, u modules.User) (int64, error)     
	DeleteUserByID(id int) (int64, error)                
}

type Repositories struct {
	UserRepository
}