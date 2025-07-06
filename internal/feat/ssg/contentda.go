package ssg

import (
	"database/sql"

	"github.com/google/uuid"
)

type ContentDA struct {
	ID        uuid.UUID      `db:"id"`
	ShortID   sql.NullString `db:"short_id"`
	UserID    sql.NullString `db:"user_id"`
	Heading   sql.NullString `db:"heading"`
	Body      sql.NullString `db:"body"`
	Status    sql.NullString `db:"status"`
	CreatedBy sql.NullString `db:"created_by"`
	UpdatedBy sql.NullString `db:"updated_by"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}
