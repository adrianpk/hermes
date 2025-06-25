package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	var users []auth.UserDA
	err = repo.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		repo.Log().Infof("User: %+v", user)
	}

	return auth.ToUsers(users), nil
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

	var user auth.UserDA
	err = repo.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return auth.User{}, err
	}

	return auth.ToUser(user), nil
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

	var userDA auth.UserExtDA
	var roles []uuid.UUID
	var permissions []uuid.UUID
	userMap := make(map[uuid.UUID]auth.User)

	for rows.Next() {
		if err := rows.StructScan(&userDA); err != nil {
			return auth.User{}, err
		}

		if _, exists := userMap[userDA.ID]; !exists {
			userMap[userDA.ID] = auth.ToUserExt(userDA)
		}

		if userDA.RoleID.Valid {
			roleID, err := uuid.Parse(userDA.RoleID.String)
			if err == nil {
				roles = append(roles, roleID)
			}
		}

		if userDA.PermissionID.Valid {
			permissionID, err := uuid.Parse(userDA.PermissionID.String)
			if err == nil {
				permissions = append(permissions, permissionID)
			}
		}
	}

	user := userMap[userDA.ID]
	user.RoleIDs = roles
	user.PermissionIDs = permissions

	return user, nil
}

func (repo *HermesRepo) CreateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "Create")
	if err != nil {
		return err
	}

	userDA := auth.ToUserDA(user)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		userDA.ID,
		userDA.Username,
		userDA.EmailEnc,
		userDA.Name,
		userDA.PasswordEnc,
		userDA.ShortID,
		userDA.CreatedBy,
		userDA.UpdatedBy,
		userDA.CreatedAt,
		userDA.UpdatedAt,
	)
	return err
}

func (repo *HermesRepo) UpdateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "Update")
	if err != nil {
		return err
	}

	userDA := auth.ToUserDA(user)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userDA.Username, userDA.EmailEnc, userDA.Name,
		userDA.ShortID, userDA.UpdatedBy, userDA.UpdatedAt, userDA.ID)
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

	userDA := auth.ToUserDA(user)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userDA.PasswordEnc, userDA.UpdatedBy, userDA.UpdatedAt, userDA.ID)
	return err
}

func (repo *HermesRepo) GetAllRoles(ctx context.Context) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resRole, "GetAll")
	if err != nil {
		return nil, err
	}

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query)
	if err != nil {
		return nil, err
	}
	return auth.ToRoles(rolesDA), nil
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

	var roleDA auth.RoleDA
	err = repo.db.GetContext(ctx, &roleDA, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Role{}, errors.New("role not found")
		}
		return auth.Role{}, err
	}
	return auth.ToRole(roleDA), nil
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

	var roleDA auth.RoleExtDA
	roleMap := make(map[uuid.UUID]auth.Role)

	for rows.Next() {
		if err := rows.StructScan(&roleDA); err != nil {
			return auth.Role{}, err
		}

		role, exists := roleMap[roleDA.ID]
		if !exists {
			role = auth.ToRoleExt(roleDA)
		}

		if roleDA.PermissionID.Valid {
			permissionID, err := uuid.Parse(roleDA.PermissionID.String)
			if err == nil {
				role.PermissionIDs = append(role.PermissionIDs, permissionID)
			}
		}

		roleMap[roleDA.ID] = role
	}

	return roleMap[roleDA.ID], nil
}

func (repo *HermesRepo) CreateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query().Get(featAuth, resRole, "Create")
	if err != nil {
		return err
	}

	roleDA := auth.ToRoleDA(role)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		roleDA.ID,
		roleDA.Name,
		roleDA.Description,
		roleDA.ShortID,
		roleDA.CreatedBy,
		roleDA.UpdatedBy,
		roleDA.CreatedAt,
		roleDA.UpdatedAt)
	return err
}

func (repo *HermesRepo) UpdateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query().Get(featAuth, resRole, "Update")
	if err != nil {
		return err
	}

	roleDA := auth.ToRoleDA(role)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		roleDA.Name,
		roleDA.Description,
		roleDA.ShortID,
		roleDA.UpdatedBy,
		roleDA.UpdatedAt,
		roleDA.ID)
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

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query)
	if err != nil {
		return nil, err
	}
	return auth.ToPermissions(permissionsDA), nil
}

// GetPermission returns a permission by ID
func (repo *HermesRepo) GetPermission(ctx context.Context, id uuid.UUID) (auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resPerm, "Get")
	if err != nil {
		return auth.Permission{}, err
	}

	var permissionDA auth.PermissionDA
	if err := repo.db.GetContext(ctx, &permissionDA, query, id); err != nil {
		if err == sql.ErrNoRows {
			return auth.Permission{}, auth.ErrPermissionNotFound
		}
		return auth.Permission{}, err
	}

	return auth.ToPermission(permissionDA), nil
}

func (repo *HermesRepo) CreatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Create")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		permissionDA.ID,
		permissionDA.ShortID,
		permissionDA.Name,
		permissionDA.Description,
		permissionDA.CreatedBy,
		permissionDA.UpdatedBy,
		permissionDA.CreatedAt,
		permissionDA.UpdatedAt,
	)
	return err
}

func (repo *HermesRepo) UpdatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Update")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		permissionDA.ShortID,
		permissionDA.Name,
		permissionDA.Description,
		permissionDA.UpdatedBy,
		permissionDA.UpdatedAt,
		permissionDA.ID)
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

	var resourcesDA []auth.ResourceDA
	err = repo.db.SelectContext(ctx, &resourcesDA, query)
	if err != nil {
		return nil, err
	}
	return auth.ToResources(resourcesDA), nil
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

	var resourceDA auth.ResourceDA
	if err := repo.db.GetContext(ctx, &resourceDA, query, id); err != nil {
		if err == sql.ErrNoRows {
			return auth.Resource{}, auth.ErrResourceNotFound
		}
		return auth.Resource{}, err
	}

	return auth.ToResource(resourceDA), nil
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

	var resourceDA auth.ResourceExtDA
	resourceMap := make(map[uuid.UUID]auth.Resource)

	for rows.Next() {
		if err := rows.StructScan(&resourceDA); err != nil {
			return auth.Resource{}, err
		}

		resource, exists := resourceMap[resourceDA.ID]
		if !exists {
			resource = auth.ToResourceExt(resourceDA)
		}

		if resourceDA.PermissionID.Valid {
			permissionID, err := uuid.Parse(resourceDA.PermissionID.String)
			if err == nil {
				resource.PermissionIDs = append(resource.PermissionIDs, permissionID)
			}
		}

		resourceMap[resourceDA.ID] = resource
	}

	return resourceMap[resourceDA.ID], nil
}

func (repo *HermesRepo) CreateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query().Get(featAuth, resRes, "Create")
	if err != nil {
		return err
	}

	resourceDA := auth.ToResourceDA(resource)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		resourceDA.ID,
		resourceDA.Name,
		resourceDA.Description,
		resourceDA.ShortID,
		resourceDA.CreatedBy,
		resourceDA.UpdatedBy,
		resourceDA.CreatedAt,
		resourceDA.UpdatedAt,
	)
	return err
}

func (repo *HermesRepo) UpdateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query().Get(featAuth, resRes, "Update")
	if err != nil {
		return err
	}

	resourceDA := auth.ToResourceDA(resource)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		resourceDA.Name,
		resourceDA.Description,
		resourceDA.ShortID,
		resourceDA.UpdatedBy,
		resourceDA.UpdatedAt,
		resourceDA.ID)
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

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query,
		userID.String(), contextType, contextID,
	)
	if err != nil {
		return nil, err
	}
	return auth.ToRoles(rolesDA), nil
}

// GetUserAssignedPermissions retrieves all permissions assigned to a user, both directly and through roles.
func (repo *HermesRepo) GetUserAssignedPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserAssignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

func (repo *HermesRepo) GetUserIndirectPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserIndirectPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetUserDirectPermissions retrieves permissions directly assigned to a user.
func (repo *HermesRepo) GetUserDirectPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserDirectPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetUserUnassignedPermissions retrieves permissions not assigned to a user, either directly or through roles.
func (repo *HermesRepo) GetUserUnassignedPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

func (repo *HermesRepo) AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resUserPerm, "AddPermissionToUser")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userID, permissionDA.ID)
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

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query,
		userID.String(), contextType, contextID,
	)
	if err != nil {
		return nil, err
	}

	return auth.ToRoles(rolesDA), nil
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

	var roleDA auth.RoleDA
	err = repo.db.GetContext(ctx, &roleDA, query, userID, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Role{}, errors.New("role not found")
		}
		return auth.Role{}, err
	}
	return auth.ToRole(roleDA), nil
}

// AddPermissionToRole adds a permission to a role.
func (repo *HermesRepo) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permission auth.Permission) error {
	query := `
		INSERT INTO role_permission (role_id, permission_id)
		VALUES (?, ?)
	`
	exec := repo.getExec(ctx)
	_, err := exec.ExecContext(ctx, query, roleID, permission.ID())
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

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, resourceID, permissionDA.ID)
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

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, roleID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetResourcePermissions returns all permissions assigned to a resource
func (repo *HermesRepo) GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resResPerm, "GetResourcePermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, resourceID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetResourceUnassignedPermissions returns all permissions not assigned to a resource
func (repo *HermesRepo) GetResourceUnassignedPermissions(ctx context.Context, resourceID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resResPerm, "GetResourceUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, resourceID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetRoleUnassignedPermissions returns all permissions not assigned to a role
func (repo *HermesRepo) GetRoleUnassignedPermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resRolePerm, "GetRoleUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, roleID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

func (r *HermesRepo) CreateOrg(ctx context.Context, org auth.Org) error {
	da := auth.ToOrgDA(org)
	query, err := r.Query().Get(featAuth, resOrg, "Create")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		da.ID, da.ShortID, da.Name, da.ShortDescription, da.Description, da.CreatedBy, da.UpdatedBy, da.CreatedAt, da.UpdatedAt,
	)
	return err
}

func (r *HermesRepo) GetDefaultOrg(ctx context.Context) (auth.Org, error) {
	query, err := r.Query().Get(featAuth, resOrg, "GetDefault")
	if err != nil {
		return auth.Org{}, err
	}
	row := r.db.QueryRowContext(ctx, query)
	var orgDA auth.OrgDA
	err = row.Scan(
		&orgDA.ID,
		&orgDA.ShortID,
		&orgDA.Name,
		&orgDA.ShortDescription,
		&orgDA.Description,
		&orgDA.CreatedBy,
		&orgDA.UpdatedBy,
		&orgDA.CreatedAt,
		&orgDA.UpdatedAt,
	)
	if err != nil {
		return auth.Org{}, err
	}
	return auth.ToOrg(orgDA), nil
}

func (r *HermesRepo) GetOrgOwners(ctx context.Context, orgID uuid.UUID) ([]auth.User, error) {
	query, err := r.Query().Get(featAuth, resOrgOwner, "GetOrgOwners")
	if err != nil {
		return nil, err
	}
	var usersDA []auth.UserDA
	err = r.db.SelectContext(ctx, &usersDA, query, orgID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
}

func (r *HermesRepo) GetOrgUnassignedOwners(ctx context.Context, orgID uuid.UUID) ([]auth.User, error) {
	query, err := r.Query().Get(featAuth, resOrgOwner, "GetOrgUnassignedOwners")
	if err != nil {
		return nil, err
	}
	var usersDA []auth.UserDA
	err = r.db.SelectContext(ctx, &usersDA, query, orgID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
}

func (r *HermesRepo) GetAllTeams(ctx context.Context, orgID uuid.UUID) ([]auth.Team, error) {
	query, err := r.Query().Get(featAuth, resTeam, "GetAll")
	if err != nil {
		return nil, err
	}
	var teamsDA []auth.TeamDA
	err = r.db.SelectContext(ctx, &teamsDA, query, orgID.String())
	if err != nil {
		return nil, err
	}
	teams := make([]auth.Team, len(teamsDA))
	for i, da := range teamsDA {
		teams[i] = auth.ToTeam(da)
	}
	return teams, nil
}

func (r *HermesRepo) GetTeam(ctx context.Context, id uuid.UUID) (auth.Team, error) {
	query, err := r.Query().Get(featAuth, resTeam, "Get")
	if err != nil {
		return auth.Team{}, err
	}
	var da auth.TeamDA
	err = r.db.GetContext(ctx, &da, query, id.String())
	if err != nil {
		return auth.Team{}, err
	}
	return auth.ToTeam(da), nil
}

func (r *HermesRepo) CreateTeam(ctx context.Context, team auth.Team) error {
	da := auth.TeamDA{
		ID:               sql.NullString{String: team.ID().String(), Valid: team.ID() != uuid.Nil},
		OrgID:            sql.NullString{String: team.OrgID.String(), Valid: team.OrgID != uuid.Nil},
		ShortID:          team.ShortID(),
		Name:             team.Name,
		ShortDescription: team.ShortDescription,
		Description:      team.Description,
		CreatedBy:        sql.NullString{String: team.CreatedBy().String(), Valid: team.CreatedBy() != uuid.Nil},
		UpdatedBy:        sql.NullString{String: team.UpdatedBy().String(), Valid: team.UpdatedBy() != uuid.Nil},
		CreatedAt:        team.CreatedAt(),
		UpdatedAt:        team.UpdatedAt(),
	}
	query, err := r.Query().Get(featAuth, resTeam, "Create")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		da.ID, da.OrgID, da.ShortID, da.Name, da.ShortDescription, da.Description, da.CreatedBy, da.UpdatedBy, da.CreatedAt, da.UpdatedAt,
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
		team.ShortID(), team.Name, team.ShortDescription, team.Description, team.UpdatedBy().String(), team.UpdatedAt(), team.ID().String(),
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
	var usersDA []auth.UserDA
	err = repo.db.SelectContext(ctx, &usersDA, query, teamID.String())
	repo.Log().Debugf("GetTeamMembers teamID: %s", teamID.String())
	for _, user := range usersDA {
		repo.Log().Debugf("User: %+v", user)
	}
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
}

func (repo *HermesRepo) GetTeamUnassignedUsers(ctx context.Context, teamID uuid.UUID) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resTeamMember, "ListUsersNotInTeam")
	if err != nil {
		return nil, err
	}
	var usersDA []auth.UserDA
	err = repo.db.SelectContext(ctx, &usersDA, query, teamID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
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

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query,
		userID.String(), "team", teamID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToRoles(rolesDA), nil
}

func (repo *HermesRepo) GetUserContextualUnassignedRoles(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetContextualUnassignedRoles")
	if err != nil {
		return nil, err
	}

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query,
		userID.String(), "team", teamID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToRoles(rolesDA), nil
}
