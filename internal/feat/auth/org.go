package auth

import (
	"encoding/json"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	orgEntityType = "org"
)

type Org struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"`
	RefValue string `json:"ref"`

	Name             string    `json:"name" db:"name"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Description      string    `json:"description" db:"description"`
	OwnerID          uuid.UUID `json:"owner_id" db:"owner_id"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

func NewOrg(name, shortDescription, description string, ownerID uuid.UUID) Org {
	o := Org{
		mType:            orgEntityType,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
		OwnerID:          ownerID,
	}
	return o
}

// Type returns the type of the entity.
func (o Org) Type() string {
	return am.DefaultType(o.mType)
}

// SetType sets the type of the entity.
func (o *Org) SetType(t string) {
	o.mType = t
}

// GetID returns the unique identifier of the entity.
func (o Org) GetID() uuid.UUID {
	return o.ID
}

// GenID delegates to the functional helper.
func (o *Org) GenID() {
	am.GenID(o)
}

// SetID sets the unique identifier of the entity.
func (o *Org) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if o.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		o.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (o *Org) GetShortID() string {
	return o.ShortID
}

// GenShortID delegates to the functional helper.
func (o *Org) GenShortID() {
	am.GenShortID(o)
}

// SetShortID sets the short ID of the entity.
func (o *Org) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if o.ShortID == "" || shouldForce {
		o.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (o *Org) TypeID() string {
	return am.Normalize(o.Type()) + "-" + o.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (o *Org) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(o, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (o *Org) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(o, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (o *Org) GetCreatedBy() uuid.UUID {
	return o.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (o *Org) GetUpdatedBy() uuid.UUID {
	return o.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (o *Org) GetCreatedAt() time.Time {
	return o.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (o *Org) GetUpdatedAt() time.Time {
	return o.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (o *Org) SetCreatedAt(t time.Time) {
	o.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (o *Org) SetUpdatedAt(t time.Time) {
	o.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (o *Org) SetCreatedBy(id uuid.UUID) {
	o.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (o *Org) SetUpdatedBy(id uuid.UUID) {
	o.UpdatedBy = id
}

// IsZero returns true if the Org is uninitialized.
func (o *Org) IsZero() bool {
	return o.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (o *Org) Slug() string {
	return am.Normalize(o.Name) + "-" + o.GetShortID()
}

// Ref returns the reference string for this entity.
func (o *Org) Ref() string {
	return o.RefValue
}

// SetRef sets the reference string for this entity.
func (o *Org) SetRef(ref string) {
	o.RefValue = ref
}

// StringID returns the unique identifier of the entity as a string.
func (o Org) StringID() string {
	return o.GetID().String()
}

// OptLabel returns the label for select options.
func (o Org) OptLabel() string {
	return o.Name
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (o *Org) UnmarshalJSON(data []byte) error {
	type Alias Org
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if o.mType == "" {
		o.mType = orgEntityType
	}

	return nil
}
