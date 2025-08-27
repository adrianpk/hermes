package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

// Resource related API handlers

func (h *APIHandler) GetAllResources(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllResources", h.Name())

	resources, err := h.service.GetAllResources(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resResourceName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resResourceNameCap)
	h.OK(w, msg, resources)
}

func (h *APIHandler) GetResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetResource", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resResourceNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	resource, err := h.service.GetResource(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resResourceName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resResourceNameCap)
	h.OK(w, msg, resource)
}

func (h *APIHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateResource", h.Name())

	var resource Resource
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newResource := NewResource(resource.Name, resource.Description, resource.Kind)
	newResource.GenCreateValues()

	err = h.service.CreateResource(r.Context(), newResource)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resResourceName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resResourceNameCap)
	h.Created(w, msg, newResource)
}

func (h *APIHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateResource", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resResourceNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var resource Resource
	err = json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedResource := NewResource(resource.Name, resource.Description, resource.Kind)
	updatedResource.SetID(id)

	err = h.service.UpdateResource(r.Context(), updatedResource)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resResourceName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resResourceNameCap)
	h.OK(w, msg, updatedResource)
}

func (h *APIHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteResource", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resResourceNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteResource(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resResourceName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resResourceNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}
