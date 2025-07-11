package auth

import (
	"database/sql"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

// ToUserDA converts a User business object to a UserDA data access object
func ToUserDA(user User) UserDA {
	return UserDA{
		ID:            user.ID(),
		ShortID:       sql.NullString{String: user.ShortID(), Valid: user.ShortID() != ""},
		Name:          sql.NullString{String: user.Name, Valid: user.Name != ""},
		Username:      sql.NullString{String: user.Username, Valid: user.Username != ""},
		EmailEnc:      user.EmailEnc,
		PasswordEnc:   user.PasswordEnc,
		RoleIDs:       toRoleIDs(user.Roles),
		PermissionIDs: toPermissionIDs(user.Permissions),
		CreatedBy:     sql.NullString{String: user.CreatedBy().String(), Valid: user.CreatedBy() != uuid.Nil},
		UpdatedBy:     sql.NullString{String: user.UpdatedBy().String(), Valid: user.UpdatedBy() != uuid.Nil},
		CreatedAt:     sql.NullTime{Time: user.CreatedAt(), Valid: !user.CreatedAt().IsZero()},
		UpdatedAt:     sql.NullTime{Time: user.UpdatedAt(), Valid: !user.UpdatedAt().IsZero()},
		LastLoginAt:   sql.NullTime{Time: derefTime(user.LastLoginAt), Valid: user.LastLoginAt != nil},
		LastLoginIP:   sql.NullString{String: user.LastLoginIP, Valid: user.LastLoginIP != ""},
		IsActive:      sql.NullBool{Bool: user.IsActive, Valid: true},
	}
}

// ToUser converts a UserDA data access object to a User business object
func ToUser(da UserDA) User {
	var lastLoginAt *time.Time
	if da.LastLoginAt.Valid {
		t := da.LastLoginAt.Time
		lastLoginAt = &t
	}

	return User{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID.String), // Added this line
			am.WithType(userType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Username:    da.Username.String,
		EmailEnc:    da.EmailEnc,
		PasswordEnc: da.PasswordEnc,
		LastLoginAt: lastLoginAt,
		LastLoginIP: da.LastLoginIP.String,
		IsActive:    da.IsActive.Bool,
	}
}

// ToUsers converts a slice of UserDA to a slice of User business objects
func ToUsers(das []UserDA) []User {
	users := make([]User, len(das))
	for i, da := range das {
		users[i] = ToUser(da)
	}
	return users
}

// ToUserExt converts UserExtDA to User including roles and permissions
func ToUserExt(da UserExtDA) User {
	user := ToUser(UserDA{
		ID:          da.ID,
		ShortID:     da.ShortID,
		Name:        da.Name,
		Username:    da.Username,
		EmailEnc:    da.EmailEnc,
		PasswordEnc: da.PasswordEnc,
		CreatedBy:   da.CreatedBy,
		UpdatedBy:   da.UpdatedBy,
		CreatedAt:   da.CreatedAt,
		UpdatedAt:   da.UpdatedAt,
		LastLoginAt: da.LastLoginAt,
		LastLoginIP: da.LastLoginIP,
		IsActive:    da.IsActive,
	})

	// Add role if present
	if da.RoleID.Valid {
		user.RoleIDs = append(user.RoleIDs, am.ParseUUIDNull(da.RoleID))
		user.Roles = append(user.Roles, Role{
			BaseModel: am.NewModel(
				am.WithID(am.ParseUUIDNull(da.RoleID)),
				am.WithType(roleType),
			),
			Name: da.RoleName.String,
		})
	}

	// Add permission if present
	if da.PermissionID.Valid {
		user.PermissionIDs = append(user.PermissionIDs, am.ParseUUIDNull(da.PermissionID))
		user.Permissions = append(user.Permissions, Permission{
			BaseModel: am.NewModel(
				am.WithID(am.ParseUUIDNull(da.PermissionID)),
				am.WithType(permissionType),
			),
			Name: da.PermissionName.String,
		})
	}

	return user
}

// ToRoleDA converts a Role business object to a RoleDA data access object
func ToRoleDA(role Role) RoleDA {
	return RoleDA{
		ID:          role.ID(),
		Name:        sql.NullString{String: role.Name, Valid: role.Name != ""},
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		ShortID:     sql.NullString{String: role.ShortID(), Valid: role.ShortID() != ""},
		Status:      sql.NullString{String: role.Status, Valid: role.Status != ""},
		Permissions: toPermissionIDs(role.Permissions),
		CreatedBy:   sql.NullString{String: role.CreatedBy().String(), Valid: role.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: role.UpdatedBy().String(), Valid: role.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: role.CreatedAt(), Valid: !role.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: role.UpdatedAt(), Valid: !role.UpdatedAt().IsZero()},
	}
}

// ToRole converts a RoleDA data access object to a Role business object
func ToRole(da RoleDA) Role {
	return Role{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID.String), // Added
			am.WithType(roleType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:          da.Name.String,
		Description:   da.Description.String,
		Status:        da.Status.String,
		PermissionIDs: da.Permissions,
		Permissions:   []Permission{},
	}
}

// ToRoles converts a slice of RoleDA to a slice of Role business objects
func ToRoles(das []RoleDA) []Role {
	roles := make([]Role, len(das))
	for i, da := range das {
		roles[i] = ToRole(da)
	}
	return roles
}

// ToRoleExt converts RoleExtDA to Role including permissions
func ToRoleExt(da RoleExtDA) Role {
	permission := Permission{
		BaseModel: am.NewModel(
			am.WithID(am.ParseUUIDNull(da.PermissionID)),
			am.WithType(permissionType),
		),
		Name: da.PermissionName.String,
	}

	return Role{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithType(roleType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
		Status:      "active", // Default status since it's not in RoleExtDA
		Permissions: []Permission{permission},
	}
}

// ToPermissionDA converts a Permission business object to a PermissionDA data access object
func ToPermissionDA(permission Permission) PermissionDA {
	return PermissionDA{
		ID:          permission.ID(),
		Name:        sql.NullString{String: permission.Name, Valid: permission.Name != ""},
		Description: sql.NullString{String: permission.Description, Valid: permission.Description != ""},
		ShortID:     sql.NullString{String: permission.ShortID(), Valid: permission.ShortID() != ""},
		CreatedBy:   sql.NullString{String: permission.CreatedBy().String(), Valid: permission.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: permission.UpdatedBy().String(), Valid: permission.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: permission.CreatedAt(), Valid: !permission.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: permission.UpdatedAt(), Valid: !permission.UpdatedAt().IsZero()},
	}
}

// ToPermission converts a PermissionDA data access object to a Permission business object
func ToPermission(da PermissionDA) Permission {
	return Permission{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID.String), // Added
			am.WithType(permissionType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
	}
}

// ToPermissions converts a slice of PermissionDA to a slice of Permission business objects
func ToPermissions(das []PermissionDA) []Permission {
	permissions := make([]Permission, len(das))
	for i, da := range das {
		permissions[i] = ToPermission(da)
	}
	return permissions
}

// ToResourceDA converts a Resource business object to a ResourceDA data access object
func ToResourceDA(resource Resource) ResourceDA {
	return ResourceDA{
		ID:          resource.ID(),
		Name:        sql.NullString{String: resource.Name, Valid: resource.Name != ""},
		Description: sql.NullString{String: resource.Description, Valid: resource.Description != ""},
		Label:       sql.NullString{String: resource.Label, Valid: resource.Label != ""},
		Type:        sql.NullString{String: resource.ResourceType, Valid: resource.ResourceType != ""},
		URI:         sql.NullString{String: resource.URI, Valid: resource.URI != ""},
		ShortID:     sql.NullString{String: resource.ShortID(), Valid: resource.ShortID() != ""},
		Permissions: toPermissionIDs(resource.Permissions),
		CreatedBy:   sql.NullString{String: resource.CreatedBy().String(), Valid: resource.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: resource.UpdatedBy().String(), Valid: resource.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: resource.CreatedAt(), Valid: !resource.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: resource.UpdatedAt(), Valid: !resource.UpdatedAt().IsZero()},
	}
}

// ToResource converts a ResourceDA data access object to a Resource business object
func ToResource(da ResourceDA) Resource {
	return Resource{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID.String), // Added
			am.WithType(resourceEntityType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:          da.Name.String,
		Description:   da.Description.String,
		Label:         da.Label.String,
		ResourceType:  da.Type.String,
		URI:           da.URI.String,
		PermissionIDs: da.Permissions,
		Permissions:   []Permission{},
	}
}

// ToResources converts a slice of ResourceDA to a slice of Resource business objects
func ToResources(das []ResourceDA) []Resource {
	resources := make([]Resource, len(das))
	for i, da := range das {
		resources[i] = ToResource(da)
	}
	return resources
}

// ToResourceExt converts ResourceExtDA to Resource including permissions
func ToResourceExt(da ResourceExtDA) Resource {
	permission := Permission{
		BaseModel: am.NewModel(
			am.WithID(am.ParseUUIDNull(da.PermissionID)),
			am.WithType(permissionType),
		),
		Name: da.PermissionName.String,
	}

	return Resource{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithType(resourceEntityType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:          da.Name.String,
		Description:   da.Description.String,
		ResourceType:  "entity",
		PermissionIDs: []uuid.UUID{am.ParseUUIDNull(da.PermissionID)},
		Permissions:   []Permission{permission},
	}
}

// OrgDA is the data access struct for Org, used for DB operations.
type OrgDA struct {
	ID               sql.NullString `db:"id"`
	ShortID          string         `db:"short_id"`
	Name             string         `db:"name"`
	ShortDescription string         `db:"short_description"`
	Description      string         `db:"description"`
	CreatedBy        sql.NullString `db:"created_by"`
	UpdatedBy        sql.NullString `db:"updated_by"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}

// ToOrg converts OrgDA to Org domain model.
func ToOrg(da OrgDA) Org {
	model := am.NewModel(
		am.WithID(am.ParseUUIDNull(da.ID)),
		am.WithShortID(da.ShortID), // Added (da.ShortID is string)
		am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
		am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
		am.WithCreatedAt(da.CreatedAt),
		am.WithUpdatedAt(da.UpdatedAt),
		am.WithType(orgEntityType),
	)
	return Org{
		BaseModel:        model,
		Name:             da.Name,
		ShortDescription: da.ShortDescription,
		Description:      da.Description,
	}
}

// ToOrgDA converts Org domain model to OrgDA for DB operations.
func ToOrgDA(org Org) OrgDA {
	return OrgDA{
		ID:               sql.NullString{String: org.ID().String(), Valid: org.ID() != uuid.Nil},
		ShortID:          org.ShortID(),
		Name:             org.Name,
		ShortDescription: org.ShortDescription,
		Description:      org.Description,
		CreatedBy:        sql.NullString{String: org.CreatedBy().String(), Valid: org.CreatedBy() != uuid.Nil},
		UpdatedBy:        sql.NullString{String: org.UpdatedBy().String(), Valid: org.UpdatedBy() != uuid.Nil},
		CreatedAt:        org.CreatedAt(),
		UpdatedAt:        org.UpdatedAt(),
	}
}

func ToTeam(da TeamDA) Team {
	model := am.NewModel(
		am.WithID(am.ParseUUIDNull(da.ID)),
		am.WithShortID(da.ShortID), // Corrected: da.ShortID is string, not sql.NullString
		am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
		am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
		am.WithCreatedAt(da.CreatedAt),
		am.WithUpdatedAt(da.UpdatedAt),
		am.WithType(teamEntityType),
	)
	return Team{
		BaseModel:        model,
		OrgID:            am.ParseUUIDNull(da.OrgID),
		Name:             da.Name,
		ShortDescription: da.ShortDescription,
		Description:      da.Description,
	}
}

func toRoleIDs(roles []Role) []uuid.UUID {
	var ids []uuid.UUID
	for _, r := range roles {
		ids = append(ids, r.ID())
	}
	return ids
}

func toPermissionIDs(perms []Permission) []uuid.UUID {
	var ids []uuid.UUID
	for _, p := range perms {
		ids = append(ids, p.ID())
	}
	return ids
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
