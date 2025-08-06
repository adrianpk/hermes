package ssg

import (
	"time"

	"github.com/google/uuid"
)

type ContentDA struct {
	ID        uuid.UUID  `db:"id"`
	ShortID   string     `db:"short_id"`
	UserID    uuid.UUID  `db:"user_id"`
	Heading   string     `db:"heading"`
	Body      string     `db:"body"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
