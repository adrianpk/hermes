package ssg

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

func (h *BFF) NewLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New layout form")
	form := NewLayoutForm(r)
	h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "", http.StatusOK)
}

func (h *BFF) CreateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create layout")

	form, err := LayoutFormFromRequest(r)
	if err != nil {
		h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil || form.HasErrors() {
		h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "Validation failed", http.StatusBadRequest)
		return
	}

	layout := ToLayout(form)
	user := h.sampleUserInSession(r)
	layout.SetCreatedBy(user.GetID())

	var createdLayout Layout
	err = h.apiClient.Post("/layouts", layout, &createdLayout)
	if err != nil {
		h.Err(w, err, "Failed to create layout via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Layout created")
	h.Redir(w, r, am.ListPath(ssgPath, layoutPath), http.StatusSeeOther)
}

func (h *BFF) EditLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit layout")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing layout ID", http.StatusBadRequest)
		return
	}

	var layout Layout
	path := fmt.Sprintf("/layouts/%s", idStr)
	err := h.apiClient.Get(path, &layout)
	if err != nil {
		h.Err(w, err, "Failed to get layout from API", http.StatusInternalServerError)
		return
	}

	form := ToLayoutForm(r, layout)
	h.renderLayoutForm(w, r, form, layout, "", http.StatusOK)
}

func (h *BFF) UpdateLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update layout")

	form, err := LayoutFormFromRequest(r)
	if err != nil {
		h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil || form.HasErrors() {
		h.renderLayoutForm(w, r, form, Newlayout("", "", ""), "Validation failed", http.StatusBadRequest)
		return
	}

	layout := ToLayout(form)
	user := h.sampleUserInSession(r)
	layout.SetUpdatedBy(user.GetID())

	path := fmt.Sprintf("/layouts/%s", layout.GetID())
	err = h.apiClient.Put(path, layout, nil)
	if err != nil {
		h.Err(w, err, "Failed to update layout via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Layout updated successfully")
	h.Redir(w, r, am.ListPath(ssgPath, layoutPath), http.StatusSeeOther)
}

func (h *BFF) ListLayouts(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List layouts")

	var layouts []Layout
	err := h.apiClient.Get("/layouts", &layouts)
	if err != nil {
		h.Err(w, err, "Failed to get layouts from API", http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, layouts)
	page.Form.SetAction(ssgPath)
	menu := page.NewMenu(ssgPath)
	menu.AddNewItem(layoutPath)

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-layouts")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *BFF) DeleteLayout(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete layout")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing layout ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/layouts/%s", idStr)
	err := h.apiClient.Delete(path)
	if err != nil {
		h.Err(w, err, "Failed to delete layout via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Layout deleted successfully")
	h.Redir(w, r, am.ListPath(ssgPath, layoutPath), http.StatusSeeOther)
}

func (h *BFF) renderLayoutForm(w http.ResponseWriter, r *http.Request, form LayoutForm, layout Layout, errorMessage string, statusCode int) {
	page := am.NewPage(r, layout)
	page.SetForm(&form)

	if layout.IsZero() {
		page.Name = "New Layout"
		page.IsNew = true
		page.Form.SetAction(am.CreatePath(ssgPath, layoutPath))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Layout"
		page.IsNew = false
		page.Form.SetAction(am.UpdatePath(ssgPath, layoutPath))
		page.Form.SetSubmitButtonText("Update")
	}

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