package ssg

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/adrianpk/hermes/internal/am"
)

const (
	contentType   = "content"
	contentTypePl = "contents"
)

// Content model.
type Content struct {
	// Common
	ID      uuid.UUID `json:"id" db:"id"`
	mType   string
	ShortID string `json:"-" db:"short_id"`

	// Content specific fields
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	SectionID uuid.UUID `json:"section_id" db:"section_id"`
	Heading   string    `json:"heading" db:"heading"`
	Body      string    `json:"body" db:"body"`
	Status    string    `json:"status" db:"status"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewContent creates a new Content.
func NewContent(heading, body string) Content {
	c := Content{
		mType:   contentType,
		Heading: heading,
		Body:    body,
	}

	return c
}

// Type returns the type of the entity.
func (c Content) Type() string {
	return am.DefaultType(c.mType)
}

// TypePlural returns the plural type of the entity.
func (c Content) TypePlural() string {
	return contentTypePl
}

// SetType sets the type of the entity.
func (c *Content) SetType(t string) {
	c.mType = t
}

// GetID returns the unique identifier of the entity.
func (c Content) GetID() uuid.UUID {
	return c.ID
}

// GenID delegates to the functional helper.
func (c *Content) GenID() {
	am.GenID(c)
}

// SetID sets the unique identifier of the entity.
func (c *Content) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if c.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		c.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (c *Content) GetShortID() string {
	return c.ShortID
}

// GenShortID delegates to the functional helper.
func (c *Content) GenShortID() {
	am.GenShortID(c)
}

// SetShortID sets the short ID of the entity.
func (c *Content) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if c.ShortID == "" || shouldForce {
		c.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (c *Content) TypeID() string {
	return am.Normalize(c.Type()) + "-" + c.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (c *Content) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(c, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (c *Content) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(c, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (c *Content) GetCreatedBy() uuid.UUID {
	return c.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (c *Content) GetUpdatedBy() uuid.UUID {
	return c.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (c *Content) GetCreatedAt() time.Time {
	return c.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (c *Content) GetUpdatedAt() time.Time {
	return c.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (c *Content) SetCreatedAt(t time.Time) {
	c.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (c *Content) SetUpdatedAt(t time.Time) {
	c.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (c *Content) SetCreatedBy(u uuid.UUID) {
	c.CreatedBy = u
}

// SetUpdatedBy implements the Auditable interface.
func (c *Content) SetUpdatedBy(u uuid.UUID) {
	c.UpdatedBy = u
}

// IsZero returns true if the Content is uninitialized.
func (c *Content) IsZero() bool {
	return c.ID == uuid.Nil
}

// Slug returns the slug for the content.
func (c Content) Slug() string {
	return "content"
}

func (c Content) OptValue() string {
	return c.GetID().String()
}

func (c Content) OptLabel() string {
	return c.Heading
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (c *Content) UnmarshalJSON(data []byte) error {
	type Alias Content
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if c.mType == "" {
		c.mType = contentType
	}

	return nil
}
