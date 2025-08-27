package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)


const (
	resContentName    = "content"
	resContentNameCap = "Content"
	resSectionName    = "section"
	resSectionNameCap = "Section"
	resLayoutName     = "layout"
	resLayoutNameCap  = "Layout"
)

type APIHandler struct {
	*am.APIHandler
	service Service
}

func NewAPIHandler(name string, service Service, options ...am.Option) *APIHandler {
	return &APIHandler{
		APIHandler: am.NewAPIHandler(name, options...),
		service:    service,
	}
}

// Layout related API handlers

func (h *APIHandler) TestLayoutsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola Mundo, Papi!")
	return
}

func (h *APIHandler) GetLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetLayout", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resLayoutNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	layout, err := h.service.GetLayout(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resLayoutNameCap)
	h.OK(w, msg, layout)
}

func (h *APIHandler) CreateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateLayout", h.Name())

	var layout Layout
	err := json.NewDecoder(r.Body).Decode(&layout)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newLayout := Newlayout(layout.Name, layout.Description, layout.Code)

	err = h.service.CreateLayout(r.Context(), newLayout)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resLayoutNameCap)
	h.Created(w, msg, newLayout)
}

func (h *APIHandler) UpdateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateLayout", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resLayoutNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var layout Layout
	err = json.NewDecoder(r.Body).Decode(&layout)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedLayout := Newlayout(layout.Name, layout.Description, layout.Code)
	updatedLayout.SetID(id, true)

	err = h.service.UpdateLayout(r.Context(), updatedLayout)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resLayoutNameCap)
	h.OK(w, msg, updatedLayout)
}

func (h *APIHandler) DeleteLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteLayout", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resLayoutNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteLayout(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resLayoutNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

func (h *APIHandler) GetAllLayouts(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllLayouts", h.Name())

	layouts, err := h.service.GetAllLayouts(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resLayoutName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resLayoutNameCap)
	h.OK(w, msg, layouts)
}

// Section related API handlers

func (h *APIHandler) GetAllSections(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllSections", h.Name())

	sections, err := h.service.GetSections(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resSectionNameCap)
	h.OK(w, msg, sections)
}

func (h *APIHandler) GetSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetSection", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resSectionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	section, err := h.service.GetSection(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resSectionNameCap)
	h.OK(w, msg, section)
}

func (h *APIHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateSection", h.Name())

	var section Section
	err := json.NewDecoder(r.Body).Decode(&section)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	newSection := NewSection(section.Name, section.Description, section.Path, section.LayoutID)

	err = h.service.CreateSection(r.Context(), newSection)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resSectionNameCap)
	h.Created(w, msg, newSection)
}

func (h *APIHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateSection", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resSectionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var section Section
	err = json.NewDecoder(r.Body).Decode(&section)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedSection := NewSection(section.Name, section.Description, section.Path, section.LayoutID)
	updatedSection.SetID(id, true)

	err = h.service.UpdateSection(r.Context(), updatedSection)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resSectionNameCap)
	h.OK(w, msg, updatedSection)
}

func (h *APIHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteSection", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resSectionNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteSection(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resSectionName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resSectionNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}

// Content related API handlers

func (h *APIHandler) GetAllContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetAllContent", h.Name())

	contents, err := h.service.GetAllContent(r.Context())
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResources, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetAllItems, resContentNameCap)
	h.OK(w, msg, contents)
}

func (h *APIHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling GetContent", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	content, err := h.service.GetContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotGetResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgGetItem, resContentNameCap)
	h.OK(w, msg, content)
}

func (h *APIHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling CreateContent", h.Name())

	var content Content
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	content.GenCreateValues()

	err = h.service.CreateContent(r.Context(), content)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotCreateResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgCreateItem, resContentNameCap)
	h.Created(w, msg, content)
}

func (h *APIHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling UpdateContent", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	var content Content
	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		h.Err(w, http.StatusBadRequest, am.ErrInvalidBody, err)
		return
	}

	updatedContent := NewContent(content.Heading, content.Body)
	updatedContent.SetID(id, true)

	err = h.service.UpdateContent(r.Context(), updatedContent)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotUpdateResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgUpdateItem, resContentNameCap)
	h.OK(w, msg, updatedContent)
}

func (h *APIHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Debugf("%s: Handling DeleteContent", h.Name())

	id, err := h.ID(w, r)
	if err != nil {
		msg := fmt.Sprintf(am.ErrInvalidID, resContentNameCap)
		h.Err(w, http.StatusBadRequest, msg, err)
		return
	}

	err = h.service.DeleteContent(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf(am.ErrCannotDeleteResource, resContentName)
		h.Err(w, http.StatusInternalServerError, msg, err)
		return
	}

	msg := fmt.Sprintf(am.MsgDeleteItem, resContentNameCap)
	h.OK(w, msg, json.RawMessage("null"))
}