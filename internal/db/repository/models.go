// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Role string

const (
	RoleStudent Role = "student"
	RoleAdmin   Role = "admin"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role `json:"role"`
	Valid bool `json:"valid"` // Valid is true if Role is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type User struct {
	ID        int64              `db:"id" json:"id"`
	Email     string             `db:"email" json:"email"`
	Password  string             `db:"password" json:"password"`
	Role      Role               `db:"role" json:"role"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
}
