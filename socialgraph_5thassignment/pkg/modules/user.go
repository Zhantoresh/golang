package modules

import "time"

type User struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Gender    string     `db:"gender" json:"gender"`
	BirthDate time.Time  `db:"birth_date" json:"birth_date"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}