package ssg

import (
	"bytes"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

const (
	contentPath    = "content"
	contentPathFmt = "%s/%s-content%s"
)

const (
	ActionNewContent    = "new-content"
	ActionCreateContent = "create-content"
	TextContent         = "Content"
)

func (h *WebHandler) NewContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New content form")
	form := NewContentForm(r)
	h.newContent(w, r, form, "", http.StatusOK)
}

func (h *WebHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create content")
	ctx := r.Context()

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.newContent(w, r, form, "Invalid form data", http.StatusBadRequest)
		return
	}

	err = form.Validate()
	if err != nil || form.HasErrors() {
		h.newContent(w, r, form, "Validation failed", http.StatusBadRequest)
		return
	}

	content := ToContentFromForm(form)
	content.GenCreateValues()

	user := h.sampleUserInSession(r)
	content.UserID = user.ID()
	content.GenCreateValues()

	err = h.service.CreateContent(ctx, content)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Content created")

	path := am.ListPath(ssgPath, contentPath)
	path = "new-content"

	h.Redir(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) newContent(w http.ResponseWriter, r *http.Request, form ContentForm, errorMessage string, statusCode int) {
	content := NewContent(form.Heading, form.Body)

	ctx := r.Context()
	sections, err := h.service.GetSections(ctx) // NOTE: This value should be cached
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, content)
	page.SetForm(form)
	page.Form.SetAction(am.CreatePath(ssgPath, contentPath))
	page.Form.SetSubmitButtonText("Create")

	page.AddSelect("sections", am.ToSelectOpt(sections))

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(content)

	tmpl, err := h.Tmpl().Get(ssgFeat, "new-content")
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

func (h *WebHandler) ListContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List content")
	ctx := r.Context()

	contents, err := h.service.GetAllContent(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, contents)
	page.Form.SetAction(ssgPath)

	menu := page.NewMenu(ssgPath)
	menu.AddNewItem(contentPath)

	tmpl, err := h.Tmpl().Get(ssgFeat, "list-content")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
