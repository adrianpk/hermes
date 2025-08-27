package auth

import (
	"github.com/adrianpk/hermes/internal/am"
)

func NewAPIRouter(handler *APIHandler, mw []am.Middleware, opts ...am.Option) *am.Router {
	core := am.NewAPIRouter("api-router", opts...)
	core.SetMiddlewares(mw)

	// User API routes
	core.Get("/users", handler.GetAllUsers)
	core.Get("/users/{id}", handler.GetUser)
	core.Post("/users", handler.CreateUser)
	core.Put("/users/{id}", handler.UpdateUser)
	core.Delete("/users/{id}", handler.DeleteUser)

	// User roles
	core.Get("/users/{userId}/roles", handler.GetUserRoles)
	core.Get("/users/{userId}/roles/unassigned", handler.GetUserUnassignedRoles)
	core.Post("/users/{userId}/roles", handler.AddRoleToUser)
	core.Delete("/users/{userId}/roles/{roleId}", handler.RemoveRoleFromUser)

	// Role API routes
	core.Get("/roles", handler.GetAllRoles)
	core.Get("/roles/{id}", handler.GetRole)
	core.Post("/roles", handler.CreateRole)
	core.Put("/roles/{id}", handler.UpdateRole)
	core.Delete("/roles/{id}", handler.DeleteRole)

	// Permission API routes
	core.Get("/permissions", handler.GetAllPermissions)
	core.Get("/permissions/{id}", handler.GetPermission)
	core.Post("/permissions", handler.CreatePermission)
	core.Put("/permissions/{id}", handler.UpdatePermission)
	core.Delete("/permissions/{id}", handler.DeletePermission)

	// Resource API routes
	core.Get("/resources", handler.GetAllResources)
	core.Get("/resources/{id}", handler.GetResource)
	core.Post("/resources", handler.CreateResource)
	core.Put("/resources/{id}", handler.UpdateResource)
	core.Delete("/resources/{id}", handler.DeleteResource)

	// Org API routes
	core.Get("/orgs", handler.GetAllOrgs)
	core.Get("/orgs/{id}", handler.GetOrg)
	core.Post("/orgs", handler.CreateOrg)
	core.Put("/orgs/{id}", handler.UpdateOrg)
	core.Delete("/orgs/{id}", handler.DeleteOrg)

	// Team API routes
	core.Get("/orgs/{orgId}/teams", handler.GetAllTeams)
	core.Post("/orgs/{orgId}/teams", handler.CreateTeam)
	core.Get("/teams/{teamId}", handler.GetTeam)
	core.Put("/teams/{teamId}", handler.UpdateTeam)
	core.Delete("/teams/{teamId}", handler.DeleteTeam)

	// Team Membership routes
	core.Get("/teams/{teamId}/members", handler.ListTeamMembers)
	core.Post("/teams/{teamId}/members", handler.AddTeamMember)
	core.Delete("/teams/{teamId}/members/{userId}", handler.RemoveTeamMember)

	// Team Member Role routes
	core.Get("/teams/{teamId}/members/{userId}/roles", handler.ListTeamMemberRoles)
	core.Post("/teams/{teamId}/members/{userId}/roles", handler.AddTeamMemberRole)
	core.Delete("/teams/{teamId}/members/{userId}/roles/{roleId}", handler.RemoveTeamMemberRole)

	// Role permissions
	core.Get("/roles/{roleId}/permissions", handler.GetRolePermissions)
	core.Get("/roles/{roleId}/permissions/unassigned", handler.GetRoleUnassignedPermissions)
	core.Post("/roles/{roleId}/permissions", handler.AddPermissionToRole)
	core.Delete("/roles/{roleId}/permissions/{permissionId}", handler.RemovePermissionFromRole)

	return core
}
