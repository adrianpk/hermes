package ssg

import (
	"database/sql"

	"github.com/google/uuid"
)

type SectionDA struct {
	ID          uuid.UUID      `db:"id"`
	ShortID     sql.NullString `db:"short_id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Path        sql.NullString `db:"path"`
	LayoutID    sql.NullString `db:"layout_id"`
	Image       sql.NullString `db:"image"`
	Header      sql.NullString `db:"header"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}
