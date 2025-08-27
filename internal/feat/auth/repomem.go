package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

// BaseRepo provides an in-memory implementation of the Repo interface.
// This implementation is intended to simplify initial prototyping.
// In the future, a relational database implementation and possibly a NoSQL implementation will be provided.
type BaseRepo struct {
	*am.BaseRepo
	mu                  sync.Mutex
	users               map[uuid.UUID]User
	roles               map[uuid.UUID]Role
	permissions         map[uuid.UUID]Permission
	resources           map[uuid.UUID]Resource
	userRoles           map[uuid.UUID][]uuid.UUID
	userPermissions     map[uuid.UUID][]uuid.UUID
	rolePermissions     map[uuid.UUID][]uuid.UUID
	resourcePermissions map[uuid.UUID][]uuid.UUID
	order               []uuid.UUID
	emailKey            []byte
}

func NewRepo(qm *am.QueryManager, opts ...am.Option) *BaseRepo {
	repo := &BaseRepo{
		BaseRepo:            am.NewRepo("todo-repo", qm, opts...),
		users:               make(map[uuid.UUID]User),
		roles:               make(map[uuid.UUID]Role),
		permissions:         make(map[uuid.UUID]Permission),
		resources:           make(map[uuid.UUID]Resource),
		userRoles:           make(map[uuid.UUID][]uuid.UUID),
		userPermissions:     make(map[uuid.UUID][]uuid.UUID),
		rolePermissions:     make(map[uuid.UUID][]uuid.UUID),
		resourcePermissions: make(map[uuid.UUID][]uuid.UUID),
		order:               []uuid.UUID{},
		emailKey:            []byte{},
	}

	repo.addSampleData() // NOTE: Used for testing purposes only.

	return repo
}

// User methods

func (repo *BaseRepo) GetAllUsers(ctx context.Context) ([]User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result []User
	for _, id := range repo.order {
		result = append(result, repo.users[id])
	}
	return result, nil
}

func (repo *BaseRepo) GetUser(ctx context.Context, id uuid.UUID, preload ...bool) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if len(preload) > 0 && preload[0] {
		return repo.getUserPreload(ctx, id)
	}
	return repo.getUser(ctx, id)
}

func (repo *BaseRepo) getUser(ctx context.Context, id uuid.UUID) (User, error) {
	user, exists := repo.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (repo *BaseRepo) getUserPreload(ctx context.Context, id uuid.UUID) (User, error) {
	user, exists := repo.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}

	user.Roles = repo.getUserRolesByID(id)
	user.Permissions = repo.getUserPermissionsByID(id)
	return user, nil
}

func (repo *BaseRepo) getUserRolesByID(userID uuid.UUID) []Role {
	var roles []Role
	for _, roleID := range repo.userRoles[userID] {
		roles = append(roles, repo.roles[roleID])
	}
	return roles
}

func (repo *BaseRepo) getUserPermissionsByID(userID uuid.UUID) []Permission {
	var permissions []Permission
	for _, permissionID := range repo.userPermissions[userID] {
		permissions = append(permissions, repo.permissions[permissionID])
	}
	return permissions
}

func (repo *BaseRepo) CreateUser(ctx context.Context, u User) (User, error) {
	// Encrypt email and password
	emailEnc, err := EncryptEmail(string(u.EmailEnc), repo.emailKey)
	if err != nil {
		return User{}, err
	}

	passwordEnc, err := HashPassword(string(u.PasswordEnc))
	if err != nil {
		return User{}, err
	}

	user := NewUser(u.Username, u.Name)
	user.SetEmailEnc(emailEnc)
	user.SetPasswordEnc(passwordEnc)
	user.RoleIDs = u.RoleIDs
	user.PermissionIDs = u.PermissionIDs

	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.GetID()]; exists {
		return User{}, errors.New("user already exists")
	}
	repo.users[user.GetID()] = user
	repo.order = append(repo.order, user.GetID())
	return user, nil
}

func (repo *BaseRepo) UpdateUser(ctx context.Context, user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.GetID()]; !exists {
		msg := fmt.Sprintf("user not found for ID: %s", user.GetID())
		return errors.New(msg)
	}
	repo.users[user.GetID()] = user
	return nil
}

func (repo *BaseRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(repo.users, id)
	for i, oid := range repo.order {
		if oid == id {
			repo.order = append(repo.order[:i], repo.order[i+1:]...)
			break
		}
	}
	delete(repo.userRoles, id)
	delete(repo.userPermissions, id)
	return nil
}

func (repo *BaseRepo) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var roles []Role
	for _, roleID := range repo.userRoles[userID] {
		roles = append(roles, repo.roles[roleID])
	}
	return roles, nil
}

func (repo *BaseRepo) AddRole(ctx context.Context, userID uuid.UUID, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}

	if _, exists := repo.roles[role.GetID()]; exists {
		return errors.New("role already exists")
	}
	repo.roles[role.GetID()] = role
	repo.userRoles[userID] = append(repo.userRoles[userID], role.GetID()) // Add role to user
	return nil
}

func (repo *BaseRepo) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}

	// Remove role from userRoles
	for i, rid := range repo.userRoles[userID] {
		if rid == roleID {
			repo.userRoles[userID] = append(repo.userRoles[userID][:i], repo.userRoles[userID][i+1:]...)
			break
		}
	}
	return nil
}

func (repo *BaseRepo) AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, exists := repo.users[userID]
	if !exists {
		return errors.New("user not found")
	}
	user.PermissionIDs = append(user.PermissionIDs, permission.GetID())
	repo.users[user.GetID()] = user
	repo.userPermissions[user.GetID()] = append(repo.userPermissions[user.GetID()], permission.GetID())
	return nil
}

func (repo *BaseRepo) RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, exists := repo.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	for i, pid := range user.PermissionIDs {
		if pid == permissionID {
			user.PermissionIDs = append(user.PermissionIDs[:i], user.PermissionIDs[i+1:]...)
			repo.users[user.GetID()] = user
			for j, upid := range repo.userPermissions[user.GetID()] {
				if upid == permissionID {
					repo.userPermissions[user.GetID()] = append(repo.userPermissions[user.GetID()][:j], repo.userPermissions[user.GetID()][j+1:]...)
					break
				}
			}
			return nil
		}
	}
	return errors.New("permission not found")
}

// Role methods

func (repo *BaseRepo) GetUserRole(ctx context.Context, userID, roleID uuid.UUID) (Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return Role{}, errors.New("user not found")
	}

	for _, rid := range repo.userRoles[userID] {
		if rid == roleID {
			return repo.roles[rid], nil
		}
	}
	return Role{}, errors.New("role not found")
}

func (repo *BaseRepo) GetRole(ctx context.Context, roleID uuid.UUID) (Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	role, exists := repo.roles[roleID]
	if !exists {
		return Role{}, errors.New("role not found")
	}
	return role, nil
}

func (repo *BaseRepo) CreateRole(ctx context.Context, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.roles[role.GetID()]; exists {
		return errors.New("role already exists")
	}
	repo.roles[role.GetID()] = role
	return nil
}

func (repo *BaseRepo) UpdateRole(ctx context.Context, userID uuid.UUID, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}

	if _, exists := repo.roles[role.GetID()]; !exists {
		msg := fmt.Sprintf("role not found for ID: %s", role.GetID())
		return errors.New(msg)
	}
	repo.roles[role.GetID()] = role
	return nil
}

func (repo *BaseRepo) DeleteRole(ctx context.Context, userID, roleID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}

	if _, exists := repo.roles[roleID]; !exists {
		return errors.New("role not found")
	}
	delete(repo.roles, roleID)

	for i, rid := range repo.userRoles[userID] {
		if rid == roleID {
			repo.userRoles[userID] = append(repo.userRoles[userID][:i], repo.userRoles[userID][i+1:]...)
			break
		}
	}

	delete(repo.rolePermissions, roleID)
	return nil
}

// AddPermissionToRole adds a permission to a role.
func (repo *BaseRepo) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	role, err := repo.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	permission, err := repo.GetPermission(ctx, permissionID)
	if err != nil {
		return err
	}

	role.Permissions = append(role.Permissions, permission)
	return nil
}

// RemovePermissionFromRole removes a permission from a role.
func (repo *BaseRepo) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	role, err := repo.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	for i, p := range role.Permissions {
		if p.GetID() == permissionID {
			role.Permissions = append(role.Permissions[:i], role.Permissions[i+1:]...)
			return nil
		}
	}

	return errors.New(am.ErrResourceNotFound)
}

// Permission methods

func (repo *BaseRepo) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var permissions []Permission
	for _, permission := range repo.permissions {
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (repo *BaseRepo) GetPermission(ctx context.Context, id uuid.UUID) (Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	permission, exists := repo.permissions[id]
	if !exists {
		return Permission{}, errors.New("permission not found")
	}
	return permission, nil
}

func (repo *BaseRepo) CreatePermission(ctx context.Context, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.permissions[permission.GetID()]; exists {
		return errors.New("permission already exists")
	}
	repo.permissions[permission.GetID()] = permission
	return nil
}

func (repo *BaseRepo) UpdatePermission(ctx context.Context, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.permissions[permission.GetID()]; !exists {
		return errors.New("permission not found")
	}
	repo.permissions[permission.GetID()] = permission
	return nil
}

func (repo *BaseRepo) DeletePermission(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.permissions[id]; !exists {
		return errors.New("permission not found")
	}
	delete(repo.permissions, id)
	return nil
}

// Resource methods

func (repo *BaseRepo) GetAllResources(ctx context.Context) ([]Resource, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var resources []Resource
	for _, resource := range repo.resources {
		resources = append(resources, resource)
	}
	return resources, nil
}

func (repo *BaseRepo) GetResource(ctx context.Context, id uuid.UUID) (Resource, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resource, exists := repo.resources[id]
	if !exists {
		return Resource{}, errors.New("resource not found")
	}
	return resource, nil
}

func (repo *BaseRepo) CreateResource(ctx context.Context, resource Resource) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[resource.GetID()]; exists {
		return errors.New("resource already exists")
	}
	repo.resources[resource.GetID()] = resource
	return nil
}

func (repo *BaseRepo) UpdateResource(ctx context.Context, resource Resource) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[resource.GetID()]; !exists {
		return errors.New("resource not found")
	}
	repo.resources[resource.GetID()] = resource
	return nil
}

func (repo *BaseRepo) DeleteResource(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[id]; !exists {
		return errors.New("resource not found")
	}
	delete(repo.resources, id)
	return nil
}

func (repo *BaseRepo) AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resource, exists := repo.resources[resourceID]
	if !exists {
		return errors.New("resource not found")
	}
	resource.PermissionIDs = append(resource.PermissionIDs, permission.GetID())
	repo.resources[resourceID] = resource
	repo.resourcePermissions[resourceID] = append(repo.resourcePermissions[resourceID], permission.GetID())
	return nil
}

func (repo *BaseRepo) RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resource, exists := repo.resources[resourceID]
	if !exists {
		return errors.New("resource not found")
	}

	for i, pid := range resource.Permissions {
		if pid.GetID() == permissionID {
			resource.Permissions = append(resource.Permissions[:i], resource.Permissions[i+1:]...)
			repo.resources[resource.GetID()] = resource
			for j, rpid := range repo.resourcePermissions[resource.GetID()] {
				if rpid == permissionID {
					repo.resourcePermissions[resource.GetID()] = append(repo.resourcePermissions[resource.GetID()][:j], repo.resourcePermissions[resource.GetID()][j+1:]...)
					break
				}
			}
			return nil
		}
	}
	return errors.New(am.ErrResourceNotFound)
}

func (repo *BaseRepo) GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[resourceID]; !exists {
		return nil, errors.New("resource not found")
	}

	var permissions []Permission
	for _, permissionID := range repo.resourcePermissions[resourceID] {
		permissions = append(permissions, repo.permissions[permissionID])
	}
	return permissions, nil
}

func (repo *BaseRepo) Debug() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result string
	result += fmt.Sprintf("%-10s %-36s %-36s %-20s\n", "Type", "ID", "Username", "Extra") // Adjusted headers, removed Slug
	for _, id := range repo.order {
		user, ok := repo.users[id]
		if !ok {
			continue
		}
		// Adjusted to use fields from user, removed Slug
		result += fmt.Sprintf("%-10s %-36s %-36s %-20s\n",
			"User", user.GetID().String(), user.Name, user.Username)
	}
	result = fmt.Sprintf("%s state:\n%s", repo.Name(), result)
	repo.Log().Info(result)
}

func (repo *BaseRepo) addSampleData() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Add sample users
	emailEnc, _ := EncryptEmail("john@example.com", repo.emailKey)
	passwordEnc, _ := HashPassword("password")
	user := NewUser("john", "John Doe")
	user.SetEmailEnc(emailEnc)
	user.SetPasswordEnc(passwordEnc)
	user.GenID() // Generate ID for the user
	user.RoleIDs = []uuid.UUID{repo.roles[uuid.MustParse("00000000-0000-0000-0000-000000000001")].GetID()}
	user.PermissionIDs = []uuid.UUID{repo.permissions[uuid.MustParse("00000000-0000-0000-0000-000000000001")].GetID()}
	repo.users[user.GetID()] = user
	repo.order = append(repo.order, user.GetID())

	// Add sample roles
	role := NewRole("admin", "Administrator", "Administrator role with full access")
	role.GenID() // Generate ID for the role
	role.PermissionIDs = []uuid.UUID{repo.permissions[uuid.MustParse("00000000-0000-0000-0000-000000000001")].GetID()}
	repo.roles[role.GetID()] = role
	repo.order = append(repo.order, role.GetID())

	// Add sample permissions
	perm := NewPermission("read", "Read permission")
	perm.GenID() // Generate ID for the permission
	repo.permissions[perm.GetID()] = perm
	repo.order = append(repo.order, perm.GetID())

	// Assign roles to users
	repo.userRoles[user.GetID()] = []uuid.UUID{role.GetID()}

	// Assign permissions to roles
	repo.rolePermissions[role.GetID()] = []uuid.UUID{perm.GetID()}

	// Add sample resources
	for i := 1; i <= 3; i++ {
		resource := NewResource(fmt.Sprintf("resource%d", i), fmt.Sprintf("resource%d description", i), "entity")
		resource.GenCreateValues()
		repo.resources[resource.GetID()] = resource
		repo.Log().Info("Created resource with ID: ", resource.GetID())
	}

	// Assign permissions to resources
	for resourceID := range repo.resources {
		repo.resourcePermissions[resourceID] = []uuid.UUID{perm.GetID()}
	}
}

func (repo *BaseRepo) GetAllRoles(ctx context.Context) ([]Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var roles []Role
	for _, role := range repo.roles {
		roles = append(roles, role)
	}
	return roles, nil
}