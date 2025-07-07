package ssg

import (
	"bytes"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

const (
	layoutPath    = "layout"
	layoutPathFmt = "%s/%s-layout%s"
)

const (
	ActionNewLayout    = "new-layout"
	ActionCreateLayout = "create-layout"
	TextLayout         = "Layout"
)

func (h *WebHandler) NewLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New layout form")
	form := NewLayoutForm(r)
	h.newLayout(w, r, form, "", http.StatusOK)
}

func (h *WebHandler) CreateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create layout")
	ctx := r.Context()

	form, err := LayoutFormFromRequest(r)
	if err != nil {
		h.newLayout(w, r, form, "Invalid form data", http.StatusBadRequest)
		return
	}

	err = form.Validate()
	if err != nil || form.HasErrors() {
		h.newLayout(w, r, form, "Validation failed", http.StatusBadRequest)
		return
	}

	layout := ToLayoutFromForm(form)
	layout.GenCreateValues()

	err = h.service.CreateLayout(ctx, layout)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Layout created")

	path := am.ListPath(ssgPath, layoutPath)
	path = "new-layout"

	h.Redir(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) newLayout(w http.ResponseWriter, r *http.Request, form LayoutForm, errorMessage string, statusCode int) {
	layout := ToLayoutFromForm(form)

	page := am.NewPage(r, layout)
	page.SetForm(form)
	page.Form.SetAction(am.CreatePath(ssgPath, layoutPath))
	page.Form.SetSubmitButtonText("Create")

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(layout)

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-layout")
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
