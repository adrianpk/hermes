package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

// Role related API handlers

func (h *APIHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllRoles", h.Name())

	roles, err := h.service.GetAllRoles(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resRoleName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resRoleNameCap)
	h.OK(w, msg, roles)
}

func (h *APIHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetRole", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resRoleNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	role, err := h.service.GetRole(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resRoleName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resRoleNameCap)
	h.OK(w, msg, role)
}

func (h *APIHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateRole", h.Name())

	var role Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newRole := NewRole(role.Name, role.Description, role.Status)
	newRole.GenCreateValues()

	err = h.service.CreateRole(r.Context(), newRole)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resRoleName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resRoleNameCap)
	h.Created(w, msg, newRole)
}

func (h *APIHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateRole", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resRoleNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var role Role
	err = json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedRole := NewRole(role.Name, role.Description, role.Status)
	updatedRole.SetID(id)

	err = h.service.UpdateRole(r.Context(), updatedRole)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resRoleName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resRoleNameCap)
	h.OK(w, msg, updatedRole)
}

func (h *APIHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteRole", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resRoleNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteRole(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resRoleName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resRoleNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Role permissions

func (h *APIHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetRolePermissions", h.Name())

	roleID, err := am.PathID(r, "roleId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid role ID in URL", err)
		return
	}

	permissions, err := h.service.GetRolePermissions(r.Context(), roleID)
	if err != nil {
		msg := fmt.Sprintf("cannot get permissions for role %s", roleID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Permissions for role %s", roleID)
	h.OK(w, msg, permissions)
}

func (h *APIHandler) GetRoleUnassignedPermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetRoleUnassignedPermissions", h.Name())

	roleID, err := am.PathID(r, "roleId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid role ID in URL", err)
		return
	}

	permissions, err := h.service.GetRoleUnassignedPermissions(r.Context(), roleID)
	if err != nil {
		msg := fmt.Sprintf("cannot get unassigned permissions for role %s", roleID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Unassigned permissions for role %s", roleID)
	h.OK(w, msg, permissions)
}

func (h *APIHandler) AddPermissionToRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddPermissionToRole", h.Name())

	roleID, err := am.PathID(r, "roleId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid role ID in URL", err)
		return
	}

	var payload struct {
		PermissionID string `json:"permission_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	permissionID, err := uuid.Parse(payload.PermissionID)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid permission ID in body", err)
		return
	}

	err = h.service.AddPermissionToRole(r.Context(), roleID, permissionID)
	if err != nil {
		msg := fmt.Sprintf("cannot add permission %s to role %s", permissionID, roleID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Permission %s added to role %s", permissionID, roleID)
	h.OK(w, msg, nil)
}

func (h *APIHandler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling RemovePermissionFromRole", h.Name())

	roleID, err := am.PathID(r, "roleId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid role ID in URL", err)
		return
	}

	permissionID, err := am.PathID(r, "permissionId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid permission ID in URL", err)
		return
	}

	err = h.service.RemovePermissionFromRole(r.Context(), roleID, permissionID)
	if err != nil {
		msg := fmt.Sprintf("cannot remove permission %s from role %s", permissionID, roleID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Permission %s removed from role %s", permissionID, roleID)
	h.OK(w, msg, nil)
}
