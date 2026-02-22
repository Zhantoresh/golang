package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var users []modules.User
	err := r.db.DB.SelectContext(ctx, &users, "SELECT id, name, email, age, city FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var user modules.User
	err := r.db.DB.GetContext(ctx, &user, "SELECT id, name, email, age, city FROM users WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(u modules.User) (int, error) {
	// Handle bad input (potential error cases)
	if u.Name == "" || u.Email == "" || u.City == "" || u.Age <= 0 {
		return 0, errors.New("invalid user fields")
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var id int
	err := r.db.DB.QueryRowxContext(ctx,
		`INSERT INTO users (name, email, age, city)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id`,
		u.Name, u.Email, u.Age, u.City,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) UpdateUser(id int, u modules.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx,
		`UPDATE users SET name=$1, email=$2, age=$3, city=$4 WHERE id=$5`,
		u.Name, u.Email, u.Age, u.City, id,
	)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil // handler will check rows==0
}

func (r *Repository) DeleteUserByID(id int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil // handler will check rows==0
}