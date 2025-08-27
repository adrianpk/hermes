package auth

import (
	"encoding/json"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	permissionType = "permission"
)

type Permission struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"`
	RefValue string `json:"ref"`

	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

func NewPermission(name, description string) Permission {
	return Permission{
		mType:       permissionType,
		Name:        name,
		Description: description,
	}
}

// Type returns the type of the entity.
func (p Permission) Type() string {
	return am.DefaultType(p.mType)
}

// SetType sets the type of the entity.
func (p *Permission) SetType(t string) {
	p.mType = t
}

// GetID returns the unique identifier of the entity.
func (p Permission) GetID() uuid.UUID {
	return p.ID
}

// GenID delegates to the functional helper.
func (p *Permission) GenID() {
	am.GenID(p)
}

// SetID sets the unique identifier of the entity.
func (p *Permission) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if p.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		p.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (p *Permission) GetShortID() string {
	return p.ShortID
}

// GenShortID delegates to the functional helper.
func (p *Permission) GenShortID() {
	am.GenShortID(p)
}

// SetShortID sets the short ID of the entity.
func (p *Permission) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if p.ShortID == "" || shouldForce {
		p.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (p *Permission) TypeID() string {
	return am.Normalize(p.Type()) + "-" + p.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (p *Permission) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(p, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (p *Permission) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(p, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (p *Permission) GetCreatedBy() uuid.UUID {
	return p.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (p *Permission) GetUpdatedBy() uuid.UUID {
	return p.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (p *Permission) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (p *Permission) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (p *Permission) SetCreatedAt(t time.Time) {
	p.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (p *Permission) SetUpdatedAt(t time.Time) {
	p.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (p *Permission) SetCreatedBy(id uuid.UUID) {
	p.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (p *Permission) SetUpdatedBy(id uuid.UUID) {
	p.UpdatedBy = id
}

// IsZero returns true if the Permission is uninitialized.
func (p *Permission) IsZero() bool {
	return p.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (p *Permission) Slug() string {
	return am.Normalize(p.Name) + "-" + p.GetShortID()
}

// Ref returns the reference string for this entity.
func (p *Permission) Ref() string {
	return p.RefValue
}

// SetRef sets the reference string for this entity.
func (p *Permission) SetRef(ref string) {
	p.RefValue = ref
}

// StringID returns the unique identifier of the entity as a string.
func (p Permission) StringID() string {
	return p.GetID().String()
}

// OptLabel returns the label for select options.
func (p Permission) OptLabel() string {
	return p.Name
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (p *Permission) UnmarshalJSON(data []byte) error {
	type Alias Permission
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if p.mType == "" {
		p.mType = permissionType
	}

	return nil
}
