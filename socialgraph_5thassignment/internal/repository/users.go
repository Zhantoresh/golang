package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"socialgraph_5thassignment/pkg/modules"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetPaginatedUsers(page, pageSize int, filters map[string]string, orderBy string, status string) (modules.PaginatedResponse, error) {
	offset := (page - 1) * pageSize

	allowedOrder := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"gender":     true,
		"birth_date": true,
	}

	if !allowedOrder[orderBy] {
		orderBy = "id"
	}

	whereParts := []string{}
	args := []interface{}{}
	argPos := 1

	if status == "" || status == "active" {
		whereParts = append(whereParts, "deleted_at IS NULL")
	} else if status == "deleted" {
		whereParts = append(whereParts, "deleted_at IS NOT NULL")
	}

	if v := filters["id"]; v != "" {
		whereParts = append(whereParts, fmt.Sprintf("id = $%d", argPos))
		args = append(args, v)
		argPos++
	}
	if v := filters["name"]; v != "" {
		whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+v+"%")
		argPos++
	}
	if v := filters["email"]; v != "" {
		whereParts = append(whereParts, fmt.Sprintf("email ILIKE $%d", argPos))
		args = append(args, "%"+v+"%")
		argPos++
	}
	if v := filters["gender"]; v != "" {
		whereParts = append(whereParts, fmt.Sprintf("gender = $%d", argPos))
		args = append(args, v)
		argPos++
	}
	if v := filters["birth_date"]; v != "" {
		whereParts = append(whereParts, fmt.Sprintf("birth_date = $%d", argPos))
		args = append(args, v)
		argPos++
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = " WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM users" + whereSQL
	var totalCount int
	if err := r.db.Get(&totalCount, countQuery, args...); err != nil {
		return modules.PaginatedResponse{}, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date, deleted_at
		FROM users
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, argPos, argPos+1)

	args = append(args, pageSize, offset)

	var users []modules.User
	if err := r.db.Select(&users, query, args...); err != nil {
		return modules.PaginatedResponse{}, err
	}

	return modules.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *UserRepository) GetCommonFriends(user1, user2 int) ([]modules.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date, u.deleted_at
		FROM users u
		JOIN user_friends f1 ON u.id = f1.friend_id
		JOIN user_friends f2 ON u.id = f2.friend_id
		WHERE f1.user_id = $1 
		  AND f2.user_id = $2
		  AND u.deleted_at IS NULL
		ORDER BY u.id
	`

	var users []modules.User
	if err := r.db.Select(&users, query, user1, user2); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) SoftDeleteUser(id int) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}