package auth

import (
	"time"

	"github.com/google/uuid"
)

// ResourceDA represents the data access layer for the Resource model.
type ResourceDA struct {
	ID            uuid.UUID    `db:"id"`
	ShortID       string       `db:"short_id"`
	Name          string       `db:"name"`
	Description   string       `db:"description"`
	Label         string       `db:"label"`
	Type          string       `db:"type"`
	URI           string       `db:"uri"`
	Permissions   []uuid.UUID  `db:"permissions"`
	CreatedBy     *string      `db:"created_by"`
	UpdatedBy     *string      `db:"updated_by"`
	CreatedAt     *time.Time   `db:"created_at"`
	UpdatedAt     *time.Time   `db:"updated_at"`
}

// ResourceExtDA represents the data access layer for the Resource with permissions.
type ResourceExtDA struct {
	ID             uuid.UUID    `db:"id"`
	ShortID        string       `db:"short_id"`
	Name           string       `db:"name"`
	Description    string       `db:"description"`
	Label          string       `db:"label"`
	Type           string       `db:"type"`
	URI            string       `db:"uri"`
	PermissionID   *string      `db:"permission_id"`
	PermissionName *string      `db:"permission_name"`
	CreatedBy      *string      `db:"created_by"`
	UpdatedBy      *string      `db:"updated_by"`
	CreatedAt      *time.Time   `db:"created_at"`
	UpdatedAt      *time.Time   `db:"updated_at"`
}