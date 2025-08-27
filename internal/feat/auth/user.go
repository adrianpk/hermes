package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

const (
	userType = "user"
)

type User struct {
	// Common
	ID       uuid.UUID `json:"id" db:"id"`
	mType    string
	ShortID  string `json:"-" db:"short_id"`
	RefValue string `json:"ref"`

	Username    string     `json:"username" db:"username"`
	Email       string     `json:"email" db:"email"`
	EmailEnc    []byte     `json:"-" db:"email_enc"`
	Name        string     `json:"name" db:"name"`
	Password    string     `json:"password" db:"password"`
	PasswordEnc []byte     `json:"-" db:"password_enc"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip,omitempty" db:"last_login_ip"`

	RoleIDs       []uuid.UUID `json:"-"`
	PermissionIDs []uuid.UUID `json:"-"`

	Roles       []Role       `json:"roles,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`

	// Audit
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
	UpdatedBy uuid.UUID `json:"-" db:"updated_by"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewUser creates a user with default values.
func NewUser(username, name string) User {
	u := User{
		mType:         userType,
		Username:      username,
		Name:          name,
		IsActive:      true,
		Roles:         []Role{},
		Permissions:   []Permission{},
		RoleIDs:       []uuid.UUID{},
		PermissionIDs: []uuid.UUID{},
	}
	return u
}

// NewUserSec creates a new user and encrypts email/password.
func NewUserSec(username, email, password, name string, emailKey []byte) (User, error) {
	emailEnc, err := EncryptEmail(email, emailKey)
	if err != nil {
		return User{}, err
	}

	passwordEnc, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	u := NewUser(username, name)
	u.SetEmailEnc(emailEnc)
	u.SetPasswordEnc(passwordEnc)

	return u, nil
}

// Type returns the type of the entity.
func (u User) Type() string {
	return am.DefaultType(u.mType)
}

// SetType sets the type of the entity.
func (u *User) SetType(t string) {
	u.mType = t
}

// GetID returns the unique identifier of the entity.
func (u User) GetID() uuid.UUID {
	return u.ID
}

// GenID delegates to the functional helper.
func (u *User) GenID() {
	am.GenID(u)
}

// SetID sets the unique identifier of the entity.
func (u *User) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if u.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		u.ID = id
	}
}

// ShortID returns the short ID portion of the slug.
func (u *User) GetShortID() string {
	return u.ShortID
}

// GenShortID delegates to the functional helper.
func (u *User) GenShortID() {
	am.GenShortID(u)
}

// SetShortID sets the short ID of the entity.
func (u *User) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if u.ShortID == "" || shouldForce {
		u.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (u *User) TypeID() string {
	return am.Normalize(u.Type()) + "-" + u.GetShortID()
}

// GenCreateValues delegates to the functional helper.
func (u *User) GenCreateValues(userID ...uuid.UUID) {
	am.SetCreateValues(u, userID...)
}

// GenUpdateValues delegates to the functional helper.
func (u *User) GenUpdateValues(userID ...uuid.UUID) {
	am.SetUpdateValues(u, userID...)
}

// CreatedBy returns the UUID of the user who created the entity.
func (u *User) GetCreatedBy() uuid.UUID {
	return u.CreatedBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (u *User) GetUpdatedBy() uuid.UUID {
	return u.UpdatedBy
}

// CreatedAt returns the creation time of the entity.
func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// UpdatedAt returns the last update time of the entity.
func (u *User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

// SetCreatedAt implements the Auditable interface.
func (u *User) SetCreatedAt(t time.Time) {
	u.CreatedAt = t
}

// SetUpdatedAt implements the Auditable interface.
func (u *User) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = t
}

// SetCreatedBy implements the Auditable interface.
func (u *User) SetCreatedBy(id uuid.UUID) {
	u.CreatedBy = id
}

// SetUpdatedBy implements the Auditable interface.
func (u *User) SetUpdatedBy(id uuid.UUID) {
	u.UpdatedBy = id
}

// IsZero returns true if the User is uninitialized.
func (u *User) IsZero() bool {
	return u.ID == uuid.Nil
}

// SetEmailEnc sets the encrypted email for the user.
func (u *User) SetEmailEnc(emailEnc []byte) {
	u.EmailEnc = emailEnc
}

// SetPasswordEnc sets the encrypted password for the user.
func (u *User) SetPasswordEnc(passwordEnc []byte) {
	u.PasswordEnc = passwordEnc
}

// AddRole adds a role to the user.
func (l *User) AddRole(role Role) {
	l.Roles = append(l.Roles, role)
}

// RemoveRole removes a role from the user.
func (l *User) RemoveRole(roleID uuid.UUID) {
	for i, role := range l.Roles {
		if role.GetID() == roleID {
			l.Roles = append(l.Roles[:i], l.Roles[i+1:]...)
			break
		}
	}
}

func (u *User) PrePersist(ctx context.Context) error {
	err := u.EncFields(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EncFields(ctx context.Context) error {
	if u.Password != "" && len(u.PasswordEnc) == 0 {
		enc, err := HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.PasswordEnc = enc
	}
	if u.Email != "" && len(u.EmailEnc) == 0 {
		key, ok := ctx.Value("encryptionKey").([]byte)
		if !ok || len(key) == 0 {
			return fmt.Errorf("encryption key not found in context")
		}
		emailEnc, err := EncryptEmail(u.Email, key)
		if err != nil {
			return err
		}
		u.EmailEnc = emailEnc
	}
	return nil
}

// UnmarshalJSON ensures model fields are initialized after unmarshal.
func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if u.mType == "" {
		u.mType = userType
	}

	return nil
}

func (u *User) Slug() string {
	return am.Normalize(u.Username) + "-" + u.GetShortID()
}

// OptLabel returns the label for select options.
func (u User) OptLabel() string {
	return u.Username
}

// Ref returns the reference string for this entity.
func (u *User) Ref() string {
	return u.RefValue
}

// SetRef sets the reference string for this entity.
func (u *User) SetRef(ref string) {
	u.RefValue = ref
}
