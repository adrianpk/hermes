package auth

import (
	"time"

	"github.com/google/uuid"
)

// UserDA represents the data access layer for the User model.
type UserDA struct {
	ID            uuid.UUID  `db:"id"`
	ShortID       string     `db:"short_id"`
	Name          string     `db:"name"`
	Username      string     `db:"username"`
	EmailEnc      []byte     `db:"email_enc"`
	PasswordEnc   []byte     `db:"password_enc"`
	RoleIDs       []uuid.UUID
	PermissionIDs []uuid.UUID
	CreatedBy     *string    `db:"created_by"`
	UpdatedBy     *string    `db:"updated_by"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	LastLoginAt   *time.Time `db:"last_login_at"`
	LastLoginIP   string     `db:"last_login_ip"`
	IsActive      bool       `db:"is_active"`
}

// UserExtDA represents the data access layer for the UserRolePermission.
type UserExtDA struct {
	ID             uuid.UUID  `db:"id"`
	ShortID        string     `db:"short_id"`
	Name           string     `db:"name"`
	Username       string     `db:"username"`
	EmailEnc       []byte     `db:"email_enc"`
	PasswordEnc    []byte     `db:"password_enc"`
	RoleID         *string    `db:"role_id"`
	PermissionID   *string    `db:"permission_id"`
	RoleName       *string    `db:"role_name"`
	PermissionName *string    `db:"permission_name"`
	CreatedBy      *string    `db:"created_by"`
	UpdatedBy      *string    `db:"updated_by"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
	LastLoginAt    *time.Time `db:"last_login_at"`
	LastLoginIP    string     `db:"last_login_ip"`
	IsActive       bool       `db:"is_active"`
}
