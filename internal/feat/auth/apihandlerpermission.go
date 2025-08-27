package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

// Permission related API handlers

func (h *APIHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllPermissions", h.Name())

	permissions, err := h.service.GetAllPermissions(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resPermissionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resPermissionNameCap)
	h.OK(w, msg, permissions)
}

func (h *APIHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetPermission", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resPermissionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	permission, err := h.service.GetPermission(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resPermissionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resPermissionNameCap)
	h.OK(w, msg, permission)
}

func (h *APIHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreatePermission", h.Name())

	var permission Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newPermission := NewPermission(permission.Name, permission.Description)
	newPermission.GenCreateValues()

	err = h.service.CreatePermission(r.Context(), newPermission)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resPermissionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resPermissionNameCap)
	h.Created(w, msg, newPermission)
}

func (h *APIHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdatePermission", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resPermissionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var permission Permission
	err = json.NewDecoder(r.Body).Decode(&permission)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedPermission := NewPermission(permission.Name, permission.Description)
	updatedPermission.SetID(id)

	err = h.service.UpdatePermission(r.Context(), updatedPermission)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resPermissionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resPermissionNameCap)
	h.OK(w, msg, updatedPermission)
}

func (h *APIHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeletePermission", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resPermissionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeletePermission(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resPermissionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resPermissionNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}
