package auth

import (
	"encoding/json"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	resourceEntityType = "resource"
)

type Resource struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"`
	RefValue string `json:"ref"`

	Name          string `json:"name" db:"name"`
	Description   string `json:"description" db:"description"`
	Label         string `json:"label" db:"label"`
	Kind          string `json:"kind" db:"kind"` // Type of resource (e.g., "url", "entity")
	URI           string `json:"uri" db:"uri"`
	PermissionIDs []uuid.UUID
	Permissions   []Permission

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

func NewResource(name, description, kind string) Resource {
	r := Resource{
		mType:         resourceEntityType,
		Name:          name,
		Description:   description,
		Kind:          kind,
		PermissionIDs: []uuid.UUID{},
		Permissions:   []Permission{},
	}
	return r
}

// Type returns the type of the entity.
func (r Resource) Type() string {
	return am.DefaultType(r.mType)
}

// SetType sets the type of the entity.
func (r *Resource) SetType(t string) {
	r.mType = t
}

// GetID returns the unique identifier of the entity.
func (r Resource) GetID() uuid.UUID {
	return r.ID
}

// GenID delegates to the functional helper.
func (r *Resource) GenID() {
	am.GenID(r)
}

// SetID sets the unique identifier of the entity.
func (r *Resource) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if r.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		r.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (r *Resource) GetShortID() string {
	return r.ShortID
}

// GenShortID delegates to the functional helper.
func (r *Resource) GenShortID() {
	am.GenShortID(r)
}

// SetShortID sets the short ID of the entity.
func (r *Resource) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if r.ShortID == "" || shouldForce {
		r.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (r *Resource) TypeID() string {
	return am.Normalize(r.Type()) + "-" + r.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (r *Resource) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(r, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (r *Resource) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(r, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (r *Resource) GetCreatedBy() uuid.UUID {
	return r.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (r *Resource) GetUpdatedBy() uuid.UUID {
	return r.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (r *Resource) GetCreatedAt() time.Time {
	return r.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (r *Resource) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (r *Resource) SetCreatedAt(t time.Time) {
	r.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (r *Resource) SetUpdatedAt(t time.Time) {
	r.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (r *Resource) SetCreatedBy(id uuid.UUID) {
	r.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (r *Resource) SetUpdatedBy(id uuid.UUID) {
	r.UpdatedBy = id
}

// IsZero returns true if the Resource is uninitialized.
func (r *Resource) IsZero() bool {
	return r.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (r *Resource) Slug() string {
	return am.Normalize(r.Name) + "-" + r.GetShortID()
}

// Ref returns the reference string for this entity.
func (r *Resource) Ref() string {
	return r.RefValue
}

// SetRef sets the reference string for this entity.
func (r *Resource) SetRef(ref string) {
	r.RefValue = ref
}

// StringID returns the unique identifier of the entity as a string.
func (r Resource) StringID() string {
	return r.GetID().String()
}

// OptLabel returns the label for select options.
func (r Resource) OptLabel() string {
	return r.Label
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (r *Resource) UnmarshalJSON(data []byte) error {
	type Alias Resource
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if r.mType == "" {
		r.mType = resourceEntityType
	}

	return nil
}
