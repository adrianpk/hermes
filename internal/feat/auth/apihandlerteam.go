package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Team related API handlers

func (h *APIHandler) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllTeams", h.Name())

	orgID, err := h.UUID(w, r, "orgId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resOrgNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	teams, err := h.service.GetAllTeams(r.Context(), orgID)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resTeamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resTeamNameCap)
	h.OK(w, msg, teams)
}

func (h *APIHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetTeam", h.Name())

	id, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	team, err := h.service.GetTeam(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resTeamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resTeamNameCap)
	h.OK(w, msg, team)
}

func (h *APIHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateTeam", h.Name())

	orgID, err := h.UUID(w, r, "orgId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resOrgNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var team Team
	err = json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newTeam := NewTeam(orgID, team.Name, team.ShortDescription, team.Description)
	newTeam.GenCreateValues()

	err = h.service.CreateTeam(r.Context(), newTeam)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resTeamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resTeamNameCap)
	h.Created(w, msg, newTeam)
}

func (h *APIHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateTeam", h.Name())

	id, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	// Get the existing team from the database
	existingTeam, err := h.service.GetTeam(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resTeamName)
		h.Err(w, http.StatusNotFound, msg, err)
		return
	}

	// Decode the request body into a temporary struct
	var payload Team
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	// Update the fields of the existing team
	existingTeam.Name = payload.Name
	existingTeam.ShortDescription = payload.ShortDescription
	existingTeam.Description = payload.Description
	existingTeam.GenUpdateValues() // Update audit fields

	// Save the updated team
	err = h.service.UpdateTeam(r.Context(), existingTeam)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resTeamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resTeamNameCap)
	h.OK(w, msg, existingTeam)
}

func (h *APIHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteTeam", h.Name())

	id, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteTeam(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resTeamName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resTeamNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Team Membership handlers

func (h *APIHandler) ListTeamMembers(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListTeamMembers", h.Name())

	teamID, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	members, err := h.service.GetTeamMembers(r.Context(), teamID)
	if err != nil {
		msg := fmt.Sprintf("Could not get members for team %s", teamID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Members for team %s retrieved successfully", teamID)
	h.OK(w, msg, members)
}

func (h *APIHandler) AddTeamMember(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddTeamMember", h.Name())

	teamID, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var payload struct {
		UserID uuid.UUID `json:"userId"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	// Using "member" as the default relation type, as seen in the webhandler
	err = h.service.AddUserToTeam(r.Context(), teamID, payload.UserID, "member")
	if err != nil {
		msg := fmt.Sprintf("Could not add user %s to team %s", payload.UserID, teamID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("User %s added to team %s successfully", payload.UserID, teamID)
	h.OK(w, msg, nil)
}

func (h *APIHandler) RemoveTeamMember(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling RemoveTeamMember", h.Name())

	teamID, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	userID, err := h.UUID(w, r, "userId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.RemoveUserFromTeam(r.Context(), teamID, userID)
	if err != nil {
		msg := fmt.Sprintf("Could not remove user %s from team %s", userID, teamID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("User %s removed from team %s successfully", userID, teamID)
	h.OK(w, msg, nil)
}

// Team Member Role handlers

func (h *APIHandler) ListTeamMemberRoles(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling ListTeamMemberRoles", h.Name())

	teamID, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	userID, err := h.UUID(w, r, "userId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	roles, err := h.service.GetUserContextualRoles(r.Context(), teamID, userID)
	if err != nil {
		msg := fmt.Sprintf("Could not get roles for user %s in team %s", userID, teamID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	unassignedRoles, err := h.service.GetUserContextualUnassignedRoles(r.Context(), teamID, userID)
	if err != nil {
		msg := fmt.Sprintf("Could not get unassigned roles for user %s in team %s", userID, teamID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	response := struct {
		AssignedRoles   []Role `json:"assignedRoles"`
		UnassignedRoles []Role `json:"unassignedRoles"`
	}{
		AssignedRoles:   roles,
		UnassignedRoles: unassignedRoles,
	}

	msg := fmt.Sprintf("Roles for user %s in team %s retrieved successfully", userID, teamID)
	h.OK(w, msg, response)
}

func (h *APIHandler) AddTeamMemberRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling AddTeamMemberRole", h.Name())

	teamID, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	userID, err := h.UUID(w, r, "userId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var payload struct {
		RoleID uuid.UUID `json:"roleId"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	err = h.service.AddContextualRole(r.Context(), userID, payload.RoleID, "team", teamID.String())
	if err != nil {
		msg := fmt.Sprintf("Could not add role %s to user %s in team %s", payload.RoleID, userID, teamID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Role %s added to user %s in team %s successfully", payload.RoleID, userID, teamID)
	h.OK(w, msg, nil)
}

func (h *APIHandler) RemoveTeamMemberRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling RemoveTeamMemberRole", h.Name())

	teamID, err := h.UUID(w, r, "teamId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resTeamNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	userID, err := h.UUID(w, r, "userId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resUserNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	roleID, err := h.UUID(w, r, "roleId")
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resRoleNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.RemoveContextualRole(r.Context(), userID, roleID, "team", teamID.String())
	if err != nil {
		msg := fmt.Sprintf("Could not remove role %s from user %s in team %s", roleID, userID, teamID)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf("Role %s removed from user %s in team %s successfully", roleID, userID, teamID)
	h.OK(w, msg, nil)
}

// UUID is a helper function to get a UUID from the URL parameters.
// It's a temporary solution until we can refactor the am.APIHandler to include it.
func (h *APIHandler) UUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, key)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
