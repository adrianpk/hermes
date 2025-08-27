package auth

import (
	"encoding/json"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	teamEntityType = "team"
)

type Team struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"`
	RefValue string `json:"ref"`

	OrgID            uuid.UUID `json:"org_id" db:"org_id"`
	OrgRef           string    `json:"org_ref" db:"org_ref"`
	Name             string    `json:"name" db:"name"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Description      string    `json:"description" db:"description"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

func NewTeam(orgID uuid.UUID, name, shortDescription, description string) Team {
	newTeam := Team{
		mType:            teamEntityType,
		OrgID:            orgID,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
	}
	return newTeam
}

// Type returns the type of the entity.
func (t Team) Type() string {
	return am.DefaultType(t.mType)
}

// SetType sets the type of the entity.
func (t *Team) SetType(newType string) {
	t.mType = newType
}

// GetID returns the unique identifier of the entity.
func (t Team) GetID() uuid.UUID {
	return t.ID
}

// GenID delegates to the functional helper.
func (t *Team) GenID() {
	am.GenID(t)
}

// SetID sets the unique identifier of the entity.
func (t *Team) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if t.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		t.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (t *Team) GetShortID() string {
	return t.ShortID
}

// GenShortID delegates to the functional helper.
func (t *Team) GenShortID() {
	am.GenShortID(t)
}

// SetShortID sets the short ID of the entity.
func (t *Team) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if t.ShortID == "" || shouldForce {
		t.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (t *Team) TypeID() string {
	return am.Normalize(t.Type()) + "-" + t.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (t *Team) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(t, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (t *Team) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(t, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (t *Team) GetCreatedBy() uuid.UUID {
	return t.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (t *Team) GetUpdatedBy() uuid.UUID {
	return t.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (t *Team) GetCreatedAt() time.Time {
	return t.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (t *Team) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (t *Team) SetCreatedAt(time time.Time) {
	t.CreatedAt = time
}

// SetUpdatedAt implements the Auditable interface.
func (t *Team) SetUpdatedAt(time time.Time) {
	t.UpdatedAt = time
}

// SetCreatedBy implements the Auditable interface.
func (t *Team) SetCreatedBy(id uuid.UUID) {
	t.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (t *Team) SetUpdatedBy(id uuid.UUID) {
	t.UpdatedBy = id
}

// IsZero returns true if the Team is uninitialized.
func (t *Team) IsZero() bool {
	return t.ID == uuid.Nil
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
func (t *Team) Slug() string {
	return am.Normalize(t.Name) + "-" + t.GetShortID()
}

// Ref returns the reference string for this entity.
func (t *Team) Ref() string {
	return t.RefValue
}

// SetRef sets the reference string for this entity.
func (t *Team) SetRef(ref string) {
	t.RefValue = ref
}



// OptLabel returns the label for select options.
func (t Team) OptLabel() string {
	return t.Name
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (t *Team) UnmarshalJSON(data []byte) error {
	type Alias Team
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if t.mType == "" {
		t.mType = teamEntityType
	}

	return nil
}
