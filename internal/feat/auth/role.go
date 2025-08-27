package auth

import (
	"encoding/json"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	roleType = "role"
)

// Role represents a todo user role.
type Role struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"`
	RefValue string `json:"ref"`

	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Status        string    `json:"status" db:"status"`
	PermissionIDs []uuid.UUID
	Permissions   []Permission

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewRole creates a new Role.
func NewRole(name, description, status string) Role {
	return Role{
		mType:         roleType,
		Name:          name,
		Description:   description,
		Status:        status,
		PermissionIDs: []uuid.UUID{},
		Permissions:   []Permission{},
	}
}

// Type returns the type of the entity.
func (r Role) Type() string {
	return am.DefaultType(r.mType)
}

// SetType sets the type of the entity.
func (r *Role) SetType(t string) {
	r.mType = t
}

// GetID returns the unique identifier of the entity.
func (r Role) GetID() uuid.UUID {
	return r.ID
}

// GenID delegates to the functional helper.
func (r *Role) GenID() {
	am.GenID(r)
}

// SetID sets the unique identifier of the entity.
func (r *Role) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if r.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		r.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (r *Role) GetShortID() string {
	return r.ShortID
}

// GenShortID delegates to the functional helper.
func (r *Role) GenShortID() {
	am.GenShortID(r)
}

// SetShortID sets the short ID of the entity.
func (r *Role) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if r.ShortID == "" || shouldForce {
		r.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (r *Role) TypeID() string {
	return am.Normalize(r.Type()) + "-" + r.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (r *Role) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(r, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (r *Role) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(r, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (r *Role) GetCreatedBy() uuid.UUID {
	return r.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (r *Role) GetUpdatedBy() uuid.UUID {
	return r.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (r *Role) GetCreatedAt() time.Time {
	return r.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (r *Role) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (r *Role) SetCreatedAt(t time.Time) {
	r.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (r *Role) SetUpdatedAt(t time.Time) {
	r.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (r *Role) SetCreatedBy(id uuid.UUID) {
	r.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (r *Role) SetUpdatedBy(id uuid.UUID) {
	r.UpdatedBy = id
}

// IsZero returns true if the Role is uninitialized.
func (r *Role) IsZero() bool {
	return r.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (r *Role) Slug() string {
	return am.Normalize(r.Name) + "-" + r.GetShortID()
}

// Ref returns the reference string for this entity.
func (r *Role) Ref() string {
	return r.RefValue
}

// SetRef sets the reference string for this entity.
func (r *Role) SetRef(ref string) {
	r.RefValue = ref
}



// OptLabel returns the label for select options.
func (r Role) OptLabel() string {
	return r.Name
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (r *Role) UnmarshalJSON(data []byte) error {
	type Alias Role
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if r.mType == "" {
		r.mType = roleType
	}

	return nil
}
