package ssg

import (
	"time"

	"github.com/google/uuid"
)

type LayoutDA struct {
	ID          uuid.UUID  `db:"id"`
	ShortID     string     `db:"short_id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Code        string     `db:"code"`
	CreatedBy   *string    `db:"created_by"`
	UpdatedBy   *string    `db:"updated_by"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
