package ssg

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

func (h *BFF) NewSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New section form")
	form := NewSectionForm(r)
	h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "", http.StatusOK)
}

func (h *BFF) CreateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create section")

	form, err := SectionFormFromRequest(r)
	if err != nil {
		h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "Invalid form data", http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil || form.HasErrors() {
		h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "Validation failed", http.StatusBadRequest)
		return
	}

	section := ToSection(form)
	user := h.sampleUserInSession(r)
	section.SetCreatedBy(user.GetID())

	var response struct {
		Section Section `json:"section"`
	}
	err = h.apiClient.Post("/sections", section, &response)
	if err != nil {
		h.Err(w, err, "Failed to create section via API", http.StatusInternalServerError)
		return
	}
	createdSection := response.Section

	h.FlashInfo(w, r, "Section created")
	h.Redir(w, r, am.EditPath(ssgPath, sectionPath, createdSection.GetID()), http.StatusSeeOther)
}

func (h *BFF) EditSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit section")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing section ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Section Section `json:"section"`
	}
	path := fmt.Sprintf("/sections/%s", idStr)
	err := h.apiClient.Get(path, &response)
	if err != nil {
		h.Err(w, err, "Failed to get section from API", http.StatusInternalServerError)
		return
	}
	section := response.Section

	form := ToSectionForm(r, section)
	h.renderSectionForm(w, r, form, section, "", http.StatusOK)
}

func (h *BFF) UpdateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update section")

	form, err := SectionFormFromRequest(r)
	if err != nil {
		h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "Invalid form data", http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil || form.HasErrors() {
		h.renderSectionForm(w, r, form, NewSection("", "", "", uuid.Nil), "Validation failed", http.StatusBadRequest)
		return
	}

	section := ToSection(form)
	user := h.sampleUserInSession(r)
	section.SetUpdatedBy(user.GetID())

	path := fmt.Sprintf("/sections/%s", section.GetID())
	err = h.apiClient.Put(path, section, nil)
	if err != nil {
		h.Err(w, err, "Failed to update section via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Section updated successfully")
	h.Redir(w, r, am.ListPath(ssgPath, sectionPath), http.StatusSeeOther)
}

func (h *BFF) ListSections(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List sections")

	var response struct {
		Sections []Section `json:"sections"`
	}
	err := h.apiClient.Get("/sections", &response)
	if err != nil {
		h.Err(w, err, "Failed to get sections from API", http.StatusInternalServerError)
		return
	}
	sections := response.Sections

	page := am.NewPage(r, sections)
	page.Form.SetAction(ssgPath)
	menu := page.NewMenu(ssgPath)
	menu.AddNewItem(sectionPath)

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-sections")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, http.StatusOK)
}

func (h *BFF) ShowSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show section")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing section ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Section Section `json:"section"`
	}
	path := fmt.Sprintf("/sections/%s", idStr)
	err := h.apiClient.Get(path, &response)
	if err != nil {
		h.Err(w, err, "Failed to get section from API", http.StatusInternalServerError)
		return
	}
	section := response.Section

	page := am.NewPage(r, section)
	page.Name = "Show Section"

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(section, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-section")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, http.StatusOK)
}

func (h *BFF) DeleteSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete section")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing section ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/sections/%s", idStr)
	err := h.apiClient.Delete(path)
	if err != nil {
		h.Err(w, err, "Failed to delete section via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Section deleted successfully")
	h.Redir(w, r, am.ListPath(ssgPath, sectionPath), http.StatusSeeOther)
}

func (h *BFF) renderSectionForm(w http.ResponseWriter, r *http.Request, form SectionForm, section Section, errorMessage string, statusCode int) {
	var response struct {
		Layouts []Layout `json:"layouts"`
	}
	err := h.apiClient.Get("/layouts", &response)
	if err != nil {
		h.Err(w, err, "Failed to get layouts from API", http.StatusInternalServerError)
		return
	}
	layouts := response.Layouts

	page := am.NewPage(r, section)
	page.SetForm(&form)
	page.AddSelect("layouts", am.ToSelectOpt(layouts))

	if section.IsZero() {
		page.Name = "New Section"
		page.IsNew = true
		page.Form.SetAction(am.CreatePath(ssgPath, sectionPath))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Section"
		page.IsNew = false
		page.Form.SetAction(am.UpdatePath(ssgPath, sectionPath))
		page.Form.SetSubmitButtonText("Update")
	}

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(section)

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-section")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	page.SetFlash(h.GetFlash(r))

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	h.OK(w, r, &buf, statusCode)
}
