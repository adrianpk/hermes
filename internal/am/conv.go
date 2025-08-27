package am

import (
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/google/uuid"
)

// NOTE: Prefer minimizing external dependencies. This library is used for now
// to provide robust pluralization,but may be replaced by a simpler custom
// solution in the future if only basic functionality is needed.
var pluralizer = pluralize.NewClient()

// Helper functions for pointer conversions
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func TimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func TimeVal(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func UUIDPtr(id uuid.UUID) *string {
	if id == uuid.Nil {
		return nil
	}
	s := id.String()
	return &s
}

func UUIDVal(s *string) uuid.UUID {
	if s == nil {
		return uuid.Nil
	}
	parsedUUID, err := uuid.Parse(*s)
	if err != nil {
		return uuid.Nil
	}
	return parsedUUID
}

func Pluralize(noun string) string {
	return pluralizer.Plural(noun)
}
