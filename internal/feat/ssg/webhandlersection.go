package ssg

import (
	"bytes"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

const (
	sectionPath    = "section"
	sectionPathFmt = "%s/%s-section%s"
)

const (
	ActionNewSection    = "new-section"
	ActionCreateSection = "create-section"
	TextSection         = "Section"
)

func (h *WebHandler) NewSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New section form")
	form := NewSectionForm(r)
	h.newSection(w, r, form, "", http.StatusOK)
}

func (h *WebHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create section")
	ctx := r.Context()

	form, err := SectionFormFromRequest(r)
	if err != nil {
		h.newSection(w, r, form, "Invalid form data", http.StatusBadRequest)
		return
	}

	err = form.Validate()
	if err != nil || form.HasErrors() {
		h.newSection(w, r, form, "Validation failed", http.StatusBadRequest)
		return
	}

	section := ToSectionFromForm(form)
	section.GenCreateValues()

	err = h.service.CreateSection(ctx, section)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Section created")

	path := am.ListPath(ssgPath, sectionPath)
	path = "new-section"

	h.Redir(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) newSection(w http.ResponseWriter, r *http.Request, form SectionForm, errorMessage string, statusCode int) {
	section := ToSectionFromForm(form)

	page := am.NewPage(r, section)
	page.SetForm(form)
	page.Form.SetAction(am.CreatePath(ssgPath, sectionPath))
	page.Form.SetSubmitButtonText("Create")

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
