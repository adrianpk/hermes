package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

// User related API handlers

func (h *APIHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllUsers", h.Name())

	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// Sanitize sensitive data before responding
	for i := range users {
		users[i].Email = ""
		users[i].Password = ""
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resUserNameCap)
	h.OK(w, msg, users)
}

func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetUser", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var user User
	embed := r.URL.Query().Get("embed")

	switch embed {
	case "permissions":
		user, err = h.service.GetUserWithPermissions(r.Context(), id)
	case "team_permissions":
		user, err = h.service.GetUserWithAllPermissions(r.Context(), id)
	default:
		user, err = h.service.GetUser(r.Context(), id)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("User with ID %s not found", id)
			h.Err(w, http.StatusNotFound, msg, err)
			return
		}
		msg := fmt.Sprintf(am.ErrCannotGetResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// Sanitize sensitive data before responding
	user.Email = ""
	user.Password = ""

	msg := fmt.Sprintf(am.MsgGetItem, resUserNameCap)
	h.OK(w, msg, user)
}

func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateUser", h.Name())

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	encKey := h.Cfg().ByteSliceVal(am.Key.SecEncryptionKey)
	newUser, err := FormToUser(UserForm{
		Username:     user.Username,
		Email:        user.Email,
		Name:         user.Name,
		Password:     user.Password,
		PasswordConf: user.Password, // Assuming password confirmation is the same
	}, encKey)
	if err != nil {
		h.Err(w, http.StatusInternalServerError, "cannot prepare user for creation", err)
		return
	}

	err = h.service.CreateUser(r.Context(), &newUser)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// Sanitize sensitive data before responding
	newUser.Email = ""
	newUser.Password = ""

	msg := fmt.Sprintf(am.MsgCreateItem, resUserNameCap)
	h.Created(w, msg, newUser)
}

func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateUser", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedUser := NewUser(user.Username, user.Name)
	updatedUser.SetID(id)

	// Update email if provided
	if user.Email != "" {
		emailEnc, err := h.crypto.EncryptEmail(user.Email)
		if err != nil {
			h.Err(w, http.StatusInternalServerError, "cannot encrypt email", err)
			return
		}
		updatedUser.EmailEnc = emailEnc
	}

	// Update password if provided
	if user.Password != "" {
		passwordEnc, err := HashPassword(user.Password)
		if err != nil {
			h.Err(w, http.StatusInternalServerError, "cannot hash password", err)
			return
		}
		updatedUser.PasswordEnc = passwordEnc
	}

	err = h.service.UpdateUser(r.Context(), &updatedUser)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	// Sanitize sensitive data before responding
	updatedUser.Email = ""
	updatedUser.Password = ""

	msg := fmt.Sprintf(am.MsgUpdateItem, resUserNameCap)
	h.OK(w, msg, updatedUser)
}

func (h *APIHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteUser", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resUserName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resUserNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// User roles

func (h *APIHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetUserRoles", h.Name())

	userID, err := am.PathID(r, "userId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid user ID in URL", err)
		return
	}

	roles, err := h.service.GetUserRoles(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("cannot get roles for user %s", userID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Roles for user %s", userID)
	h.OK(w, msg, roles)
}

func (h *APIHandler) GetUserUnassignedRoles(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetUserUnassignedRoles", h.Name())

	userID, err := am.PathID(r, "userId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid user ID in URL", err)
		return
	}

	roles, err := h.service.GetUserUnassignedRoles(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("cannot get unassigned roles for user %s", userID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Unassigned roles for user %s", userID)
	h.OK(w, msg, roles)
}

func (h *APIHandler) AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddRoleToUser", h.Name())

	userID, err := am.PathID(r, "userId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid user ID in URL", err)
		return
	}

	var payload struct {
		RoleID string `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	roleID, err := uuid.Parse(payload.RoleID)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid role ID in body", err)
		return
	}

	err = h.service.AddRole(r.Context(), userID, roleID)
	if err != nil {
		msg := fmt.Sprintf("cannot add role %s to user %s", roleID, userID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Role %s added to user %s", roleID, userID)
	h.OK(w, msg, nil)
}

func (h *APIHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling RemoveRoleFromUser", h.Name())

	userID, err := am.PathID(r, "userId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid user ID in URL", err)
		return
	}

	roleID, err := am.PathID(r, "roleId")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid role ID in URL", err)
		return
	}

	err = h.service.RemoveRole(r.Context(), userID, roleID)
	if err != nil {
		msg := fmt.Sprintf("cannot remove role %s from user %s", roleID, userID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Role %s removed from user %s", roleID, userID)
	h.OK(w, msg, nil)
}
