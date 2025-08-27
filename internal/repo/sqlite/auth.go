package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/adrianpk/hermes/internal/feat/auth"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	featAuth      = "auth"
	resUser       = "user"
	resRole       = "role"
	resPerm       = "permission"
	resRes        = "resource"
	resUserRole   = "user_role"
	resUserPerm   = "user_permission"
	resRolePerm   = "role_permission"
	resResPerm    = "resource_permission"
	resOrg        = "org"
	resOrgOwner   = "org_owner"
	resTeamMember = "team_member"
	resTeam       = "team"
)

func (repo *HermesRepo) GetUsers(ctx context.Context) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resUser, "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []auth.User
	for rows.Next() {
		var (
			id uuid.UUID
			username string
			emailEnc []byte
			passwordEnc []byte
			name string
			shortID string
			createdBy uuid.UUID
			updatedBy uuid.UUID
			createdAt time.Time
			updatedAt time.Time
			lastLoginAt sql.NullTime
			lastLoginIP sql.NullString
			isActive bool
		)

		err := rows.Scan(&id, &username, &emailEnc, &passwordEnc, &name, &shortID, &createdBy, &updatedBy, &createdAt, &updatedAt, &lastLoginAt, &lastLoginIP, &isActive)
		if err != nil {
			return nil, err
		}

		user := auth.NewUser(username, name)
		user.SetID(id)
		user.SetShortID(shortID)
		user.SetEmailEnc(emailEnc)
		user.SetPasswordEnc(passwordEnc)
		user.SetCreatedBy(createdBy)
		user.SetUpdatedBy(updatedBy)
		user.SetCreatedAt(createdAt)
		user.SetUpdatedAt(updatedAt)
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}
		if lastLoginIP.Valid {
			user.LastLoginIP = lastLoginIP.String
		}
		user.IsActive = isActive
		
		users = append(users, user)
	}

	return users, nil
}

func (repo *HermesRepo) GetUser(ctx context.Context, id uuid.UUID, preload ...bool) (auth.User, error) {
	if len(preload) > 0 && preload[0] {
		return repo.getUserPreload(ctx, id)
	}
	return repo.getUser(ctx, id)
}

func (repo *HermesRepo) getUser(ctx context.Context, id uuid.UUID) (auth.User, error) {
	query, err := repo.Query().Get(featAuth, resUser, "Get")
	if err != nil {
		return auth.User{}, err
	}

	row := repo.db.QueryRowxContext(ctx, query, id)

	var (
		userID uuid.UUID
		name string
		username string
		emailEnc []byte
		passwordEnc []byte
		shortID string
		createdBy uuid.UUID
		updatedBy uuid.UUID
		createdAt time.Time
		updatedAt time.Time
		lastLoginAt sql.NullTime
		lastLoginIP sql.NullString
		isActive bool
	)

	err = row.Scan(&userID, &name, &username, &emailEnc, &passwordEnc, &shortID, &createdBy, &updatedBy, &createdAt, &updatedAt, &lastLoginAt, &lastLoginIP, &isActive)
	if err != nil {
		return auth.User{}, err
	}

	user := auth.NewUser(username, name)
	user.SetID(userID)
	user.SetShortID(shortID)
	user.SetEmailEnc(emailEnc)
	user.SetPasswordEnc(passwordEnc)
	user.SetCreatedBy(createdBy)
	user.SetUpdatedBy(updatedBy)
	user.SetCreatedAt(createdAt)
	user.SetUpdatedAt(updatedAt)
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	if lastLoginIP.Valid {
		user.LastLoginIP = lastLoginIP.String
	}
	user.IsActive = isActive

	return user, nil
}

func (repo *HermesRepo) getUserPreload(ctx context.Context, id uuid.UUID) (auth.User, error) {
	query, err := repo.Query().Get(featAuth, resUser, "GetPreload")
	if err != nil {
		return auth.User{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, id)
	if err != nil {
		return auth.User{}, err
	}
	defer rows.Close()

	var user auth.User
	userMap := make(map[uuid.UUID]auth.User)
	roleMap := make(map[uuid.UUID]*auth.Role) // Use pointer to modify in place
	permissionMap := make(map[uuid.UUID]*auth.Permission) // Use pointer to modify in place

	for rows.Next() {
		var (
			userID        uuid.UUID
			userName      string
			userUsername  string
			userEmailEnc  []byte
			userPasswordEnc []byte
			userShortID   string
			userCreatedBy uuid.UUID
			userUpdatedBy uuid.UUID
			userCreatedAt time.Time
			userUpdatedAt time.Time
			userLastLoginAt sql.NullTime
			userLastLoginIP sql.NullString
			userIsActive  bool
			roleID        sql.NullString
			roleName      sql.NullString
			permissionID  sql.NullString
			permissionName sql.NullString
		)

		err := rows.Scan(
			&userID, &userName, &userUsername, &userEmailEnc, &userPasswordEnc, &userShortID,
			&userCreatedBy, &userUpdatedBy, &userCreatedAt, &userUpdatedAt, &userLastLoginAt, &userLastLoginIP, &userIsActive,
			&roleID, &roleName,
			&permissionID, &permissionName,
		)
		if err != nil {
			return auth.User{}, err
		}

		if _, exists := userMap[userID]; !exists {
			user = auth.NewUser(userUsername, userName)
			user.SetID(userID)
			user.SetShortID(userShortID)
			user.SetEmailEnc(userEmailEnc)
			user.SetPasswordEnc(userPasswordEnc)
			user.SetCreatedBy(userCreatedBy)
			user.SetUpdatedBy(userUpdatedBy)
			user.SetCreatedAt(userCreatedAt)
			user.SetUpdatedAt(userUpdatedAt)
			if userLastLoginAt.Valid {
				user.LastLoginAt = &userLastLoginAt.Time
			}
			if userLastLoginIP.Valid {
				user.LastLoginIP = userLastLoginIP.String
			}
			user.IsActive = userIsActive
			user.Roles = []auth.Role{}
			user.Permissions = []auth.Permission{}
			userMap[userID] = user
		}

		// Add role if present and not already added
		if roleID.Valid {
			parsedRoleID, err := uuid.Parse(roleID.String)
			if err != nil {
				return auth.User{}, err
			}
			if _, exists := roleMap[parsedRoleID]; !exists {
				role := auth.NewRole(roleName.String, "", "")
				role.SetID(parsedRoleID)
				roleMap[parsedRoleID] = &role
				tempUser := userMap[userID]
				tempUser.Roles = append(tempUser.Roles, role)
				userMap[userID] = tempUser
			}
		}

		// Add permission if present and not already added
		if permissionID.Valid {
			parsedPermissionID, err := uuid.Parse(permissionID.String)
			if err != nil {
				return auth.User{}, err
			}
			if _, exists := permissionMap[parsedPermissionID]; !exists {
				permission := auth.NewPermission(permissionName.String, "")
				permission.SetID(parsedPermissionID)
				permissionMap[parsedPermissionID] = &permission
				tempUser := userMap[userID]
				tempUser.Permissions = append(tempUser.Permissions, permission)
				userMap[userID] = tempUser
			}
		}
	}

	if len(userMap) == 0 {
		return auth.User{}, sql.ErrNoRows
	}

	for _, u := range userMap {
		return u, nil
	}
	return auth.User{}, sql.ErrNoRows // Should not be reached
}

func (repo *HermesRepo) CreateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "Create")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		user.GetID(),
		user.Username,
		user.EmailEnc,
		user.Name,
		user.PasswordEnc,
		user.GetShortID(),
		user.GetCreatedBy(),
		user.GetUpdatedBy(),
		user.GetCreatedAt(),
		user.GetUpdatedAt(),
	)
	return err
}

func (repo *HermesRepo) UpdateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "Update")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		user.Username,
		user.EmailEnc,
		user.Name,
		user.GetShortID(),
		user.GetUpdatedBy(),
		user.GetUpdatedAt(),
		user.GetID(),
	)
	return err
}

func (repo *HermesRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resUser, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id)
	return err
}

func (repo *HermesRepo) UpdatePassword(ctx context.Context, user auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "UpdatePassword")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		user.PasswordEnc,
		user.GetUpdatedBy(),
		user.GetUpdatedAt(),
		user.GetID(),
	)
	return err
}

func (repo *HermesRepo) GetAllRoles(ctx context.Context) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resRole, "GetAll")
	if err != nil {
		return nil, err
	}

	var roles []auth.Role
	err = repo.db.SelectContext(ctx, &roles, query)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRole retrieves a role by its ID, optionally preloading its associated permissions.
func (repo *HermesRepo) GetRole(ctx context.Context, id uuid.UUID, preload ...bool) (auth.Role, error) {
	if len(preload) > 0 && preload[0] {
		return repo.getRolePreload(ctx, id)
	}
	return repo.getRole(ctx, id)
}

func (repo *HermesRepo) getRole(ctx context.Context, id uuid.UUID) (auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resRole, "Get")
	if err != nil {
		return auth.Role{}, err
	}

	var role auth.Role
	err = repo.db.GetContext(ctx, &role, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Role{}, errors.New("role not found")
		}
		return auth.Role{}, err
	}
	return role, nil
}

func (repo *HermesRepo) getRolePreload(ctx context.Context, id uuid.UUID) (auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resRole, "GetPreload")
	if err != nil {
		return auth.Role{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, id)
	if err != nil {
		return auth.Role{}, err
	}
	defer rows.Close()

	var role auth.Role
	roleMap := make(map[uuid.UUID]auth.Role)
	permissionMap := make(map[uuid.UUID]*auth.Permission)

	for rows.Next() {
		var (
			roleID            uuid.UUID
			roleName          string
			roleDescription   string
			roleShortID       string
			roleCreatedBy     uuid.UUID
			roleUpdatedBy     uuid.UUID
			roleCreatedAt     time.Time
			roleUpdatedAt     time.Time
			permissionID      sql.NullString
			permissionName    sql.NullString
			permissionShortID sql.NullString
		)

		err := rows.Scan(
			&roleID, &roleName, &roleDescription, &roleShortID,
			&roleCreatedBy, &roleUpdatedBy, &roleCreatedAt, &roleUpdatedAt,
			&permissionID, &permissionName, &permissionShortID,
		)
		if err != nil {
			return auth.Role{}, err
		}

		if _, exists := roleMap[roleID]; !exists {
			role = auth.NewRole(roleName, roleDescription, "")
			role.SetID(roleID)
			role.SetShortID(roleShortID)
			role.SetCreatedBy(roleCreatedBy)
			role.SetUpdatedBy(roleUpdatedBy)
			role.SetCreatedAt(roleCreatedAt)
			role.SetUpdatedAt(roleUpdatedAt)
			role.Permissions = []auth.Permission{}
			roleMap[roleID] = role
		}

		if permissionID.Valid {
			parsedPermissionID, err := uuid.Parse(permissionID.String)
			if err != nil {
				return auth.Role{}, err
			}
			if _, exists := permissionMap[parsedPermissionID]; !exists {
				permission := auth.NewPermission(permissionName.String, "")
				permission.SetID(parsedPermissionID)
				if permissionShortID.Valid {
					permission.SetShortID(permissionShortID.String)
				}
				permissionMap[parsedPermissionID] = &permission
				
                // Get the role from the map, append the permission, and update the map
                tempRole := roleMap[roleID]
				tempRole.Permissions = append(tempRole.Permissions, permission)
                roleMap[roleID] = tempRole
			}
		}
	}

	if len(roleMap) == 0 {
		return auth.Role{}, sql.ErrNoRows
	}

	for _, r := range roleMap {
		return r, nil
	}
	return auth.Role{}, sql.ErrNoRows
}

func (repo *HermesRepo) CreateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query().Get(featAuth, resRole, "Create")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		role.GetID(),
		role.Name,
		role.Description,
		role.GetShortID(),
		role.GetCreatedBy(),
		role.GetUpdatedBy(),
		role.GetCreatedAt(),
		role.GetUpdatedAt(),
	)
	return err
}

func (repo *HermesRepo) UpdateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query().Get(featAuth, resRole, "Update")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		role.Name,
		role.Description,
		role.GetShortID(),
		role.GetUpdatedBy(),
		role.GetUpdatedAt(),
		role.GetID(),
	)
	return err
}

func (repo *HermesRepo) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resRole, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, roleID)
	return err
}

func (repo *HermesRepo) GetAllPermissions(ctx context.Context) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resPerm, "GetAll")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	err = repo.db.SelectContext(ctx, &permissions, query)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermission returns a permission by ID
func (repo *HermesRepo) GetPermission(ctx context.Context, id uuid.UUID) (auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resPerm, "Get")
	if err != nil {
		return auth.Permission{}, err
	}

	var permission auth.Permission
	if err := repo.db.GetContext(ctx, &permission, query, id); err != nil {
		if err == sql.ErrNoRows {
			return auth.Permission{}, auth.ErrPermissionNotFound
		}
		return auth.Permission{}, err
	}

	return permission, nil
}

func (repo *HermesRepo) CreatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Create")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		permission.GetID(),
		permission.Name,
		permission.Description,
		permission.GetShortID(),
		permission.GetCreatedBy(),
		permission.GetUpdatedBy(),
		permission.GetCreatedAt(),
		permission.GetUpdatedAt(),
	)
	return err
}

func (repo *HermesRepo) UpdatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Update")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		permission.Name,
		permission.Description,
		permission.GetShortID(),
		permission.GetUpdatedBy(),
		permission.GetUpdatedAt(),
		permission.GetID(),
	)
	return err
}

func (repo *HermesRepo) DeletePermission(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id)
	return err
}

func (repo *HermesRepo) GetAllResources(ctx context.Context) ([]auth.Resource, error) {
	query, err := repo.Query().Get(featAuth, resRes, "GetAll")
	if err != nil {
		return nil, err
	}

	var resources []auth.Resource
	err = repo.db.SelectContext(ctx, &resources, query)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// GetResource retrieves a resource by its ID, optionally preloading its associated permissions.
func (repo *HermesRepo) GetResource(ctx context.Context, id uuid.UUID, preload ...bool) (auth.Resource, error) {
	if len(preload) > 0 && preload[0] {
		return repo.getResourcePreload(ctx, id)
	}
	return repo.getResource(ctx, id)
}

func (repo *HermesRepo) getResource(ctx context.Context, id uuid.UUID) (auth.Resource, error) {
	query, err := repo.Query().Get(featAuth, resRes, "Get")
	if err != nil {
		return auth.Resource{}, err
	}

	var resource auth.Resource
	if err := repo.db.GetContext(ctx, &resource, query, id); err != nil {
		if err == sql.ErrNoRows {
			return auth.Resource{}, auth.ErrResourceNotFound
		}
		return auth.Resource{}, err
	}

	return resource, nil
}

func (repo *HermesRepo) getResourcePreload(ctx context.Context, id uuid.UUID) (auth.Resource, error) {
	query, err := repo.Query().Get(featAuth, resRes, "GetPreload")
	if err != nil {
		return auth.Resource{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, id)
	if err != nil {
		return auth.Resource{}, err
	}
	defer rows.Close()

	var resource auth.Resource
	resourceMap := make(map[uuid.UUID]auth.Resource)
	permissionMap := make(map[uuid.UUID]*auth.Permission)

	for rows.Next() {
		var (
			resourceID          uuid.UUID
			resourceName        string
			resourceDescription string
			resourceShortID     string
			resourceCreatedBy   uuid.UUID
			resourceUpdatedBy   uuid.UUID
			resourceCreatedAt   time.Time
			resourceUpdatedAt   time.Time
			permissionID        sql.NullString
			permissionName      sql.NullString
			permissionShortID   sql.NullString
		)

		err := rows.Scan(
			&resourceID, &resourceName, &resourceDescription, &resourceShortID,
			&resourceCreatedBy, &resourceUpdatedBy, &resourceCreatedAt, &resourceUpdatedAt,
			&permissionID, &permissionName, &permissionShortID,
		)
		if err != nil {
			return auth.Resource{}, err
		}

		if _, exists := resourceMap[resourceID]; !exists {
			// The NewResource function requires resourceType, but it is not in the query.
			// I will pass an empty string for now.
			resource = auth.NewResource(resourceName, resourceDescription, "")
			resource.SetID(resourceID)
			resource.SetShortID(resourceShortID)
			resource.SetCreatedBy(resourceCreatedBy)
			resource.SetUpdatedBy(resourceUpdatedBy)
			resource.SetCreatedAt(resourceCreatedAt)
			resource.SetUpdatedAt(resourceUpdatedAt)
			resource.Permissions = []auth.Permission{}
			resourceMap[resourceID] = resource
		}

		if permissionID.Valid {
			parsedPermissionID, err := uuid.Parse(permissionID.String)
			if err != nil {
				return auth.Resource{}, err
			}
			if _, exists := permissionMap[parsedPermissionID]; !exists {
				permission := auth.NewPermission(permissionName.String, "")
				permission.SetID(parsedPermissionID)
				if permissionShortID.Valid {
					permission.SetShortID(permissionShortID.String)
				}
				permissionMap[parsedPermissionID] = &permission
				
                tempResource := resourceMap[resourceID]
				tempResource.Permissions = append(tempResource.Permissions, permission)
                resourceMap[resourceID] = tempResource
			}
		}
	}

	if len(resourceMap) == 0 {
		return auth.Resource{}, sql.ErrNoRows
	}

	for _, r := range resourceMap {
		return r, nil
	}
	return auth.Resource{}, sql.ErrNoRows
}

func (repo *HermesRepo) CreateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query().Get(featAuth, resRes, "Create")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		resource.GetID(),
		resource.Name,
		resource.Description,
		resource.GetShortID(),
		resource.GetCreatedBy(),
		resource.GetUpdatedBy(),
		resource.GetCreatedAt(),
		resource.GetUpdatedAt(),
	)
	return err
}

func (repo *HermesRepo) UpdateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query().Get(featAuth, resRes, "Update")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		resource.Name,
		resource.Description,
		resource.GetShortID(),
		resource.GetUpdatedBy(),
		resource.GetUpdatedAt(),
		resource.GetID(),
	)
	return err
}

func (repo *HermesRepo) DeleteResource(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resRes, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id)
	return err
}

func (repo *HermesRepo) GetUserAssignedRoles(ctx context.Context, userID uuid.UUID, contextType, contextID string) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetUserAssignedRoles")
	if err != nil {
		return nil, err
	}

	var roles []auth.Role
	err = repo.db.SelectContext(ctx, &roles, query,
		userID.String(), contextType, contextID,
	)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetUserAssignedPermissions retrieves all permissions assigned to a user, both directly and through roles.
func (repo *HermesRepo) GetUserAssignedPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserAssignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	err = repo.db.SelectContext(ctx, &permissions, query, userID, userID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (repo *HermesRepo) GetUserIndirectPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserIndirectPermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	err = repo.db.SelectContext(ctx, &permissions, query, userID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetUserDirectPermissions retrieves permissions directly assigned to a user.
func (repo *HermesRepo) GetUserDirectPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserDirectPermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	err = repo.db.SelectContext(ctx, &permissions, query, userID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetUserUnassignedPermissions retrieves permissions not assigned to a user, either directly or through roles.
func (repo *HermesRepo) GetUserUnassignedPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	err = repo.db.SelectContext(ctx, &permissions, query, userID, userID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (repo *HermesRepo) AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resUserPerm, "AddPermissionToUser")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userID, permission.GetID())
	return err
}

func (repo *HermesRepo) RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resUserPerm, "RemovePermissionFromUser")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userID, permissionID)
	return err
}

func (repo *HermesRepo) GetUserUnassignedRoles(ctx context.Context, userID uuid.UUID, contextType, contextID string) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetUserUnassignedRoles")
	if err != nil {
		return nil, err
	}

	var roles []auth.Role
	err = repo.db.SelectContext(ctx, &roles, query,
		userID.String(), contextType, contextID,
	)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (repo *HermesRepo) AddRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, contextType, contextID string) error {
	query, err := repo.Query().Get(featAuth, resUserRole, "AddRole")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		userID.String(), roleID.String(), contextType, contextID,
		roleID.String(),
	)
	return err
}

func (repo *HermesRepo) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, contextType, contextID string) error {
	query, err := repo.Query().Get(featAuth, resUserRole, "RemoveRole")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		userID.String(), roleID.String(), contextType, contextID,
	)
	return err
}

func (repo *HermesRepo) GetUserRole(ctx context.Context, userID, roleID uuid.UUID) (auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetUserRole")
	if err != nil {
		return auth.Role{}, err
	}

	var role auth.Role
	err = repo.db.GetContext(ctx, &role, query, userID, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Role{}, errors.New("role not found")
		}
		return auth.Role{}, err
	}
	return role, nil
}

// AddPermissionToRole adds a permission to a role.
func (repo *HermesRepo) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	query := `
		INSERT INTO role_permission (role_id, permission_id)
		VALUES (?, ?)
	`
	exec := repo.getExec(ctx)
	_, err := exec.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to add permission to role: %w", err)
	}
	return nil
}

// RemovePermissionFromRole removes a permission from a role.
func (repo *HermesRepo) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	query := `
		DELETE FROM role_permission
		WHERE role_id = ? AND permission_id = ?
	`
	exec := repo.getExec(ctx)
	result, err := exec.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New(am.ErrResourceNotFound)
	}

	return nil
}

func (repo *HermesRepo) AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resResPerm, "AddPermissionToResource")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, resourceID, permission.GetID())
	return err
}

func (repo *HermesRepo) RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resResPerm, "RemovePermissionFromResource")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, resourceID, permissionID)
	return err
}

// GetRolePermissions returns all permissions assigned to a role
func (repo *HermesRepo) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resRolePerm, "GetRolePermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	if err := repo.db.SelectContext(ctx, &permissions, query, roleID); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetResourcePermissions returns all permissions assigned to a resource
func (repo *HermesRepo) GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resResPerm, "GetResourcePermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	if err := repo.db.SelectContext(ctx, &permissions, query, resourceID); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetResourceUnassignedPermissions returns all permissions not assigned to a resource
func (repo *HermesRepo) GetResourceUnassignedPermissions(ctx context.Context, resourceID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resResPerm, "GetResourceUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	if err := repo.db.SelectContext(ctx, &permissions, query, resourceID); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetRoleUnassignedPermissions returns all permissions not assigned to a role
func (repo *HermesRepo) GetRoleUnassignedPermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resRolePerm, "GetRoleUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissions []auth.Permission
	if err := repo.db.SelectContext(ctx, &permissions, query, roleID); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *HermesRepo) CreateOrg(ctx context.Context, org auth.Org) error {
	query, err := r.Query().Get(featAuth, resOrg, "Create")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		org.GetID(),
		org.Name,
		org.ShortDescription,
		org.Description,
		org.OwnerID,
		org.GetShortID(),
		org.GetCreatedBy(),
		org.GetUpdatedBy(),
		org.GetCreatedAt(),
		org.GetUpdatedAt(),
	)
	return err
}

func (r *HermesRepo) GetDefaultOrg(ctx context.Context) (auth.Org, error) {
	query, err := r.Query().Get(featAuth, resOrg, "GetDefault")
	if err != nil {
		return auth.Org{}, err
	}
	var org auth.Org
	err = r.db.GetContext(ctx, &org, query)
	if err != nil {
		return auth.Org{}, err
	}
	return org, nil
}

func (r *HermesRepo) GetOrgOwners(ctx context.Context, orgID uuid.UUID) ([]auth.User, error) {
	query, err := r.Query().Get(featAuth, resOrgOwner, "GetOrgOwners")
	if err != nil {
		return nil, err
	}
	var users []auth.User
	err = r.db.SelectContext(ctx, &users, query, orgID.String())
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *HermesRepo) GetOrgUnassignedOwners(ctx context.Context, orgID uuid.UUID) ([]auth.User, error) {
	query, err := r.Query().Get(featAuth, resOrgOwner, "GetOrgUnassignedOwners")
	if err != nil {
		return nil, err
	}
	var users []auth.User
	err = r.db.SelectContext(ctx, &users, query, orgID.String())
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *HermesRepo) GetAllTeams(ctx context.Context, orgID uuid.UUID) ([]auth.Team, error) {
	query, err := r.Query().Get(featAuth, resTeam, "GetAll")
	if err != nil {
		return nil, err
	}
	var teams []auth.Team
	err = r.db.SelectContext(ctx, &teams, query, orgID.String())
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *HermesRepo) GetTeam(ctx context.Context, id uuid.UUID) (auth.Team, error) {
	query, err := r.Query().Get(featAuth, resTeam, "Get")
	if err != nil {
		return auth.Team{}, err
	}
	var team auth.Team
	err = r.db.GetContext(ctx, &team, query, id.String())
	if err != nil {
		return auth.Team{}, err
	}
	return team, nil
}

func (r *HermesRepo) CreateTeam(ctx context.Context, team auth.Team) error {
	query, err := r.Query().Get(featAuth, resTeam, "Create")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		team.GetID(),
		team.OrgID,
		team.Name,
		team.ShortDescription,
		team.Description,
		team.GetShortID(),
		team.GetCreatedBy(),
		team.GetUpdatedBy(),
		team.GetCreatedAt(),
		team.GetUpdatedAt(),
	)
	return err
}

func (r *HermesRepo) UpdateTeam(ctx context.Context, team auth.Team) error {
	query, err := r.Query().Get(featAuth, resTeam, "Update")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		team.OrgID,
		team.Name,
		team.ShortDescription,
		team.Description,
		team.GetShortID(),
		team.GetUpdatedBy(),
		team.GetUpdatedAt(),
		team.GetID(),
	)
	return err
}

func (r *HermesRepo) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	query, err := r.Query().Get(featAuth, resTeam, "Delete")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id.String())
	return err
}

func (r *HermesRepo) AddOrgOwner(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	query, err := r.Query().Get(featAuth, resOrgOwner, "Add")
	if err != nil {
		return err
	}
	id := uuid.New()
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id.String(), orgID.String(), userID.String())
	return err
}

func (r *HermesRepo) RemoveOrgOwner(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	query, err := r.Query().Get(featAuth, resOrgOwner, "Remove")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, orgID.String(), userID.String())
	return err
}

func (repo *HermesRepo) GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resTeamMember, "ListTeamMembers")
	repo.Log().Debugf("GetTeamMembers query: %s", query)
	if err != nil {
		return nil, err
	}
	var users []auth.User
	err = repo.db.SelectContext(ctx, &users, query, teamID.String())
	repo.Log().Debugf("GetTeamMembers teamID: %s", teamID.String())
	// The loop below is for debugging purposes and can be removed once the refactoring is complete.
	for _, user := range users {
		repo.Log().Debugf("User: %+v", user)
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *HermesRepo) GetTeamUnassignedUsers(ctx context.Context, teamID uuid.UUID) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resTeamMember, "ListUsersNotInTeam")
	if err != nil {
		return nil, err
	}
	var users []auth.User
	err = repo.db.SelectContext(ctx, &users, query, teamID.String())
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *HermesRepo) AddUserToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, relationType string) error {
	query, err := repo.Query().Get(featAuth, resTeamMember, "AddUserToTeam")
	if err != nil {
		return err
	}
	id := uuid.New().String()
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id, teamID.String(), userID.String(), relationType)
	return err
}

func (repo *HermesRepo) RemoveUserFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resTeamMember, "RemoveUserFromTeam")
	if err != nil {
		return err
	}
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, teamID.String(), userID.String())
	return err
}

func (repo *HermesRepo) GetUserContextualRoles(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetContextualAssignedRoles")
	if err != nil {
		return nil, err
	}

	var roles []auth.Role
	err = repo.db.SelectContext(ctx, &roles, query,
		userID.String(), "team", teamID.String())
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (repo *HermesRepo) GetUserContextualUnassignedRoles(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetContextualUnassignedRoles")
	if err != nil {
		return nil, err
	}

	var roles []auth.Role
	err = repo.db.SelectContext(ctx, &roles, query,
		userID.String(), "team", teamID.String())
	if err != nil {
		return nil, err
	}
	return roles, nil
}
