package am

import (
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// Model interface composes Identifiable, Auditable, and Stampable interfaces.
type Model interface {
	Identifiable
	Auditable
	Stampable
	// Seedable
}

// Identifiable interface represents an entity with an ID, Slug, and TypeID.
type Identifiable interface {
	// Type returns the type of the entity.
	Type() string
	// GetID returns the unique identifier of the entity.
	GetID() uuid.UUID
	// GenID generates and sets the unique identifier of the entity if it is not set yet.
	GenID()
	// SetID sets the unique identifier of the entity.
	SetID(id uuid.UUID, force ...bool)
	// ShortID returns the short ID portion of the slug.
	GetShortID() string
	// GenShortID generates and sets the short ID if it is not set yet.
	GenShortID()
	// SetShortID sets the short ID of the entity.
	SetShortID(shortID string, force ...bool)
	// TypeID returns a universal identifier for a specific model instance.
	TypeID() string
	// Slug returns the slug of the entity,
	Slug() string
}

// Auditable interface represents an entity with audit information.
type Auditable interface {
	// CreatedBy returns the UUID of the user who created the entity.
	GetCreatedBy() uuid.UUID
	// UpdatedBy returns the UUID of the user who last updated the entity.
	GetUpdatedBy() uuid.UUID
	// CreatedAt returns the creation time of the entity.
	GetCreatedAt() time.Time
	// UpdatedAt returns the last update time of the entity.
	GetUpdatedAt() time.Time
	// Setters for audit fields
	SetCreatedBy(uuid.UUID)
	SetUpdatedBy(uuid.UUID)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}

type Stampable interface {
	GenCreateValues(userID ...uuid.UUID) // Modified
	GenUpdateValues(userID ...uuid.UUID) // Modified
}

// Seedable is an interface for entities that need special logic before being inserted into the database during seeding.
type Seedable interface {
	// Ref returns the reference string for this entity (used in seed data for relationships).
	Ref() string
	// SetRef sets the reference string for this entity.
	SetRef(ref string)
}

// --- Functional Helpers ---

// GenID generates a new UUID for a model if it doesn't have one.
func GenID(i Identifiable) {
	if i.GetID() == uuid.Nil {
		i.SetID(uuid.New(), true) // Force setting
	}
}

// GenShortID generates a new short ID for a model if it doesn't have one.
func GenShortID(i Identifiable) {
	if i.GetShortID() == "" {
		newUUID := uuid.New()
		segments := strings.Split(newUUID.String(), "-")
		i.SetShortID(segments[len(segments)-1], true) // Force setting
	}
}

// SetCreateValues sets the initial values for a new model.
func SetCreateValues(m Model, userID ...uuid.UUID) {
	GenID(m)
	GenShortID(m)
	now := time.Now()
	m.SetCreatedAt(now)
	m.SetUpdatedAt(now)
	if len(userID) > 0 && userID[0] != uuid.Nil {
		m.SetCreatedBy(userID[0])
		m.SetUpdatedBy(userID[0])
	}
}

// SetUpdateValues updates the timestamp for a model modification.
func SetUpdateValues(m Model, userID ...uuid.UUID) {
	m.SetUpdatedAt(time.Now())
	if len(userID) > 0 {
		m.SetUpdatedBy(userID[0])
	}
}

// Normalize utility function
func Normalize(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, r := range str {
		if unicode.IsSpace(r) {
			b.WriteRune('-')
		} else if r > unicode.MaxASCII {
			b.WriteRune('-')
		} else {
			b.WriteRune(r)
		}
	}
	return strings.ToLower(b.String())
}

func DefaultType(currentType string) string {
	if currentType == "" {
		return "model"
	}
	return currentType
}
