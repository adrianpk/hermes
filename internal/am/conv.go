package am

import (
	"time"

	"github.com/google/uuid"
)

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
