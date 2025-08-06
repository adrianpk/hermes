package auth

import (
	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

// ToUserDA converts a User business object to a UserDA data access object
func ToUserDA(user User) UserDA {
	return UserDA{
		ID:          user.ID(),
		ShortID:     user.ShortID(),
		Name:        user.Name,
		Username:    user.Username,
		EmailEnc:    user.EmailEnc,
		PasswordEnc: user.PasswordEnc,
		RoleIDs:       toRoleIDs(user.Roles),
		PermissionIDs: toPermissionIDs(user.Permissions),
		CreatedBy:   am.UUIDPtr(user.CreatedBy()),
		UpdatedBy:   am.UUIDPtr(user.UpdatedBy()),
		CreatedAt:   am.TimePtr(user.CreatedAt()),
		UpdatedAt:   am.TimePtr(user.UpdatedAt()),
		LastLoginAt: user.LastLoginAt,
		LastLoginIP: user.LastLoginIP,
		IsActive:    user.IsActive,
	}
}

// ToUser converts a UserDA data access object to a User business object
func ToUser(da UserDA) User {
	return User{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(userType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:        da.Name,
		Username:    da.Username,
		EmailEnc:    da.EmailEnc,
		PasswordEnc: da.PasswordEnc,
		LastLoginAt: da.LastLoginAt,
		LastLoginIP: da.LastLoginIP,
		IsActive:    da.IsActive,
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
	if da.RoleID != nil {
		user.RoleIDs = append(user.RoleIDs, am.UUIDVal(da.RoleID))
		user.Roles = append(user.Roles, Role{
			BaseModel: am.NewModel(
				am.WithID(am.UUIDVal(da.RoleID)),
				am.WithType(roleType),
			),
			Name: am.StringVal(da.RoleName),
		})
	}

	// Add permission if present
	if da.PermissionID != nil {
		user.PermissionIDs = append(user.PermissionIDs, am.UUIDVal(da.PermissionID))
		user.Permissions = append(user.Permissions, Permission{
			BaseModel: am.NewModel(
				am.WithID(am.UUIDVal(da.PermissionID)),
				am.WithType(permissionType),
			),
			Name: am.StringVal(da.PermissionName),
		})
	}

	return user
}

// ToRoleDA converts a Role business object to a RoleDA data access object
func ToRoleDA(role Role) RoleDA {
	return RoleDA{
		ID:          role.ID(),
		Name:        role.Name,
		Description: role.Description,
		ShortID:     role.ShortID(),
		Status:      role.Status,
		Permissions: toPermissionIDs(role.Permissions),
		CreatedBy:   am.UUIDPtr(role.CreatedBy()),
		UpdatedBy:   am.UUIDPtr(role.UpdatedBy()),
		CreatedAt:   am.TimePtr(role.CreatedAt()),
		UpdatedAt:   am.TimePtr(role.UpdatedAt()),
	}
}

// ToRole converts a RoleDA data access object to a Role business object
func ToRole(da RoleDA) Role {
	return Role{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(roleType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:          da.Name,
		Description:   da.Description,
		Status:        da.Status,
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
			am.WithID(am.UUIDVal(da.PermissionID)),
			am.WithType(permissionType),
		),
		Name: am.StringVal(da.PermissionName),
	}

	return Role{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithType(roleType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:        da.Name,
		Description: da.Description,
		Status:      "active", // Default status since it's not in RoleExtDA
		Permissions: []Permission{permission},
	}
}

// ToPermissionDA converts a Permission business object to a PermissionDA data access object
func ToPermissionDA(permission Permission) PermissionDA {
	return PermissionDA{
		ID:          permission.ID(),
		Name:        permission.Name,
		Description: permission.Description,
		ShortID:     permission.ShortID(),
		CreatedBy:   am.UUIDPtr(permission.CreatedBy()),
		UpdatedBy:   am.UUIDPtr(permission.UpdatedBy()),
		CreatedAt:   am.TimePtr(permission.CreatedAt()),
		UpdatedAt:   am.TimePtr(permission.UpdatedAt()),
	}
}

// ToPermission converts a PermissionDA data access object to a Permission business object
func ToPermission(da PermissionDA) Permission {
	return Permission{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(permissionType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:        da.Name,
		Description: da.Description,
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
		Name:        resource.Name,
		Description: resource.Description,
		ShortID:     resource.ShortID(),
		Label:       resource.Label,
		Type:        resource.ResourceType,
		URI:         resource.URI,
		Permissions: toPermissionIDs(resource.Permissions),
		CreatedBy:   am.UUIDPtr(resource.CreatedBy()),
		UpdatedBy:   am.UUIDPtr(resource.UpdatedBy()),
		CreatedAt:   am.TimePtr(resource.CreatedAt()),
		UpdatedAt:   am.TimePtr(resource.UpdatedAt()),
	}
}

// ToResource converts a ResourceDA data access object to a Resource business object
func ToResource(da ResourceDA) Resource {
	return Resource{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(resourceEntityType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:          da.Name,
		Description:   da.Description,
		Label:         da.Label,
		ResourceType:  da.Type,
		URI:           da.URI,
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
			am.WithID(am.UUIDVal(da.PermissionID)),
			am.WithType(permissionType),
		),
		Name: am.StringVal(da.PermissionName),
	}

	return Resource{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithType(resourceEntityType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:          da.Name,
		Description:   da.Description,
		ResourceType:  "entity",
		PermissionIDs: []uuid.UUID{am.UUIDVal(da.PermissionID)},
		Permissions:   []Permission{permission},
	}
}

// ToOrgDA converts a Org business object to a OrgDA data access object
func ToOrgDA(org Org) OrgDA {
	return OrgDA{
		ID:               org.ID(),
		ShortID:          org.ShortID(),
		Name:             org.Name,
		ShortDescription: org.ShortDescription,
		Description:      org.Description,
		CreatedBy:        am.UUIDPtr(org.CreatedBy()),
		UpdatedBy:        am.UUIDPtr(org.UpdatedBy()),
		CreatedAt:        am.TimePtr(org.CreatedAt()),
		UpdatedAt:        am.TimePtr(org.UpdatedAt()),
	}
}

// ToOrg converts a OrgDA data access object to a Org business object
func ToOrg(da OrgDA) Org {
	return Org{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(orgEntityType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:             da.Name,
		ShortDescription: da.ShortDescription,
		Description:      da.Description,
	}
}

// ToTeamDA converts a Team business object to a TeamDA data access object
func ToTeamDA(team Team) TeamDA {
	return TeamDA{
		ID:               team.ID(),
		ShortID:          team.ShortID(),
		OrgID:            team.OrgID,
		Name:             team.Name,
		ShortDescription: team.ShortDescription,
		Description:      team.Description,
		CreatedBy:        am.UUIDPtr(team.CreatedBy()),
		UpdatedBy:        am.UUIDPtr(team.UpdatedBy()),
		CreatedAt:        am.TimePtr(team.CreatedAt()),
		UpdatedAt:        am.TimePtr(team.UpdatedAt()),
	}
}

// ToTeam converts a TeamDA data access object to a Team business object
func ToTeam(da TeamDA) Team {
	return Team{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(teamEntityType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		OrgID:            da.OrgID,
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