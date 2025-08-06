package auth

import (
	"time"

	"github.com/google/uuid"
)

// OrgDA represents the data access layer for the Org model.
type OrgDA struct {
	ID               uuid.UUID  `db:"id"`
	ShortID          string     `db:"short_id"`
	Name             string     `db:"name"`
	ShortDescription string     `db:"short_description"`
	Description      string     `db:"description"`
	CreatedBy        *string    `db:"created_by"`
	UpdatedBy        *string    `db:"updated_by"`
	CreatedAt        *time.Time `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}
