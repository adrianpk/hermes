package am

import (
	"database/sql"

	"github.com/google/uuid"
)

func ParseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}

func ParseUUIDNull(s sql.NullString) uuid.UUID {
	if s.Valid {
		u, err := uuid.Parse(s.String)
		if err != nil {
			return uuid.Nil
		}
		return u
	}
	return uuid.Nil
}
