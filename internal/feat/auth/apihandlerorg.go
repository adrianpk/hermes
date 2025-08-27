package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

// Org related API handlers

func (h *APIHandler) GetAllOrgs(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllOrgs", h.Name())

	orgs, err := h.service.GetAllOrgs(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resOrgName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resOrgNameCap)
	h.OK(w, msg, orgs)
}

func (h *APIHandler) GetOrg(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetOrg", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resOrgNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	org, err := h.service.GetOrg(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resOrgName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resOrgNameCap)
	h.OK(w, msg, org)
}

func (h *APIHandler) CreateOrg(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateOrg", h.Name())

	var org Org
	err := json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newOrg := NewOrg(org.Name, org.ShortDescription, org.Description)
	newOrg.GenCreateValues()

	err = h.service.CreateOrg(r.Context(), newOrg)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resOrgName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resOrgNameCap)
	h.Created(w, msg, newOrg)
}

func (h *APIHandler) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateOrg", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resOrgNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var org Org
	err = json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedOrg := NewOrg(org.Name, org.ShortDescription, org.Description)
	updatedOrg.SetID(id)

	err = h.service.UpdateOrg(r.Context(), updatedOrg)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resOrgName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resOrgNameCap)
	h.OK(w, msg, updatedOrg)
}

func (h *APIHandler) DeleteOrg(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteOrg", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resOrgNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteOrg(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resOrgName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resOrgNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}
