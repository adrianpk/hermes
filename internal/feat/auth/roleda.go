package auth

import (
	"time"

	"github.com/google/uuid"
)

// RoleDA represents the data access layer for the Role model.
type RoleDA struct {
	ID            uuid.UUID    `db:"id"`
	ShortID       string       `db:"short_id"`
	Name          string       `db:"name"`
	Description   string       `db:"description"`
	Contextual    bool         `db:"contextual"`
	Status        string       `db:"status"`
	Permissions   []uuid.UUID  `db:"permissions"`
	CreatedBy     *string      `db:"created_by"`
	UpdatedBy     *string      `db:"updated_by"`
	CreatedAt     *time.Time   `db:"created_at"`
	UpdatedAt     *time.Time   `db:"updated_at"`
}

// RoleExtDA represents the data access layer for the Role with permissions.
type RoleExtDA struct {
	ID             uuid.UUID    `db:"id"`
	ShortID        string       `db:"short_id"`
	Name           string       `db:"name"`
	Description    string       `db:"description"`
	Contextual     bool         `db:"contextual"`
	Status         string       `db:"status"`
	PermissionID   *string      `db:"permission_id"`
	PermissionName *string      `db:"permission_name"`
	CreatedBy      *string      `db:"created_by"`
	UpdatedBy      *string      `db:"updated_by"`
	CreatedAt      *time.Time   `db:"created_at"`
	UpdatedAt      *time.Time   `db:"updated_at"`
}