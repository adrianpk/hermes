package auth

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

// Role handlers
func (h *WebHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List roles")
	ctx := r.Context()

	roles, err := h.service.GetAllRoles(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, roles)
	page.Form.SetAction(authPath)

	menu := page.NewMenu(authPath)
	menu.AddNewItem(rolePath)

	tmpl, err := h.tm.Get("auth", "list-roles")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) NewRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New role form")

	role := NewRole("", "", "active")

	page := am.NewPage(r, role)
	page.Form.SetAction(am.CreatePath(authPath, rolePath))
	page.Form.SetSubmitButtonText("Create")

	menu := page.NewMenu(authPath)
	menu.AddListItem(role)

	tmpl, err := h.tm.Get("auth", "new-role")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create role")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	description := r.Form.Get("description")
	status := r.Form.Get("status")
	if status == "" {
		status = "active"
	}

	role := NewRole(name, description, status)
	role.GenCreateValues()

	err := h.service.CreateRole(ctx, role)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, rolePath)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) ShowRole(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Show role ", id)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(r, role)

	menu := page.NewMenu(authPath)
	menu.AddListItem(role)
	menu.AddEditItem(role)
	menu.AddDeleteItem(role)
	menu.AddGenericItem("list-role-permissions", role.ID().String(), "Permissions")

	tmpl, err := h.tm.Get("auth", "show-role")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) EditRole(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Edit role ", id)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(r, role)
	page.Form.SetAction(am.UpdatePath(authPath, rolePath))
	page.Form.SetSubmitButtonText("Update")

	menu := page.NewMenu(authPath)
	menu.AddListItem(role)

	tmpl, err := h.tm.Get("auth", "edit-role")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	idStr := r.Form.Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Update role ", id)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	role.Name = r.Form.Get("name")
	role.Description = r.Form.Get("description")

	err = h.service.UpdateRole(ctx, role)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, rolePath)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Delete role ", id)
	ctx := r.Context()

	err = h.service.DeleteRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, rolePath)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

// Role relationships
func (h *WebHandler) ListRolePermissions(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	h.Log().Info("Showing role permissions ", "id ", id)

	ctx := r.Context()
	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	assigned, err := h.service.GetRolePermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	unassigned, err := h.service.GetRoleUnassignedPermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, struct {
		ID                   uuid.UUID
		Name                 string
		Description          string
		Permissions          []Permission
		AvailablePermissions []Permission
	}{
		ID:                   role.ID(),
		Name:                 role.Name,
		Description:          role.Description,
		Permissions:          assigned,
		AvailablePermissions: unassigned,
	})

	menu := page.NewMenu(authPath)
	menu.AddListItem(role)
	menu.AddEditItem(role)
	menu.AddDeleteItem(role)
	menu.AddGenericItem("list-role-permissions", role.ID().String(), "Permissions")

	tmpl, err := h.tm.Get("auth", "list-role-permissions")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) AddPermissionToRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add permission to role")
	ctx := r.Context()

	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	err = h.service.AddPermissionToRole(ctx, roleID, permissionID)
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			path := am.ListPath(authPath, listRolePermissionsPath) + "?id=" + roleID.String()
			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		}
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, listRolePermissionsPath) + "?id=" + roleID.String()
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove permission from role")
	ctx := r.Context()

	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemovePermissionFromRole(ctx, roleID, permissionID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, listRolePermissionsPath) + "?id=" + roleID.String()
	http.Redirect(w, r, path, http.StatusSeeOther)
}

// Handler methods for contextual roles
func (h *WebHandler) ListUserContextualRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamIDStr := r.URL.Query().Get("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		h.Err(w, err, "Invalid team ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	team, err := h.service.GetTeam(ctx, teamID)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	assignedRoles, err := h.service.GetUserContextualRoles(ctx, teamID, userID)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	unassignedRoles, err := h.service.GetUserContextualUnassignedRoles(ctx, teamID, userID)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, struct {
		User            User
		Team            Team
		AssignedRoles   []Role
		UnassignedRoles []Role
	}{
		User:            user,
		Team:            team,
		AssignedRoles:   assignedRoles,
		UnassignedRoles: unassignedRoles,
	})

	page.Form.SetAction(am.CreatePath(authPath, contextualRolePath))

	menu := am.NewMenu(authPath)
	menu.AddGenericItem("list-team-members", team.ID().String(), "Back")

	tmpl, err := h.tm.Get("auth", "list-user-contextual-roles")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) AddContextualRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add contextual role to user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}
	teamIDStr := r.FormValue("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	err = h.service.AddContextualRole(ctx, userID, roleID, "team", teamID.String())
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, listUserContextualRolesPath) + "?team_id=" + teamID.String() + "&user_id=" + userID.String()
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) RemoveContextualRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove contextual role from user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}
	teamIDStr := r.FormValue("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	err = h.service.RemoveContextualRole(ctx, userID, roleID, "team", teamID.String())
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, listUserContextualRolesPath) + "?team_id=" + teamID.String() + "&user_id=" + userID.String()
	http.Redirect(w, r, path, http.StatusSeeOther)
}
