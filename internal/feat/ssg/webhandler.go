package ssg

import (
	"bytes"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

const (
	ssgPath     = "/ssg"
	contentPath = "content"

	userPathFmt = "%s/%s-content%s"
)

const (
	ActionNewContent    = "new-content"
	ActionCreateContent = "create-content"
	TextContent         = "Content"
)

type WebHandler struct {
	*am.Handler
	service Service
	tm      *am.TemplateManager
	flash   *am.FlashManager
}

func NewWebHandler(tm *am.TemplateManager, flash *am.FlashManager, service Service, options ...am.Option) *WebHandler {
	handler := am.NewHandler("web-handler", options...)
	return &WebHandler{
		Handler: handler,
		service: service,
		tm:      tm,
		flash:   flash,
	}
}

func (h *WebHandler) NewContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New content form")

	form := ContentForm{
		BaseForm: am.NewBaseForm(r),
	}

	h.renderNewContent(w, r, form, "", http.StatusOK)
}

func (h *WebHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create content")
	ctx := r.Context()

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.renderNewContent(w, r, form, "Invalid form data", http.StatusBadRequest)
		return
	}

	err = form.Validate()
	if err != nil || form.HasErrors() {
		h.renderNewContent(w, r, form, "Validation failed", http.StatusBadRequest)
		return
	}

	content := NewContent(form.Heading, form.Body)
	content.GenCreateValues()

	err = h.service.CreateContent(ctx, content)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(ssgPath, contentPath)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) renderNewContent(w http.ResponseWriter, r *http.Request, form ContentForm, errorMessage string, statusCode int) {
	h.Log().Info("Rendering new content page with errors")

	content := NewContent(form.Heading, form.Body)

	page := am.NewPage[Content](r, content)
	page.SetForm(form)
	page.Form.SetAction(am.CreatePath(ssgPath, contentPath))
	page.Form.SetSubmitButtonText("Create")

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(content)

	tmpl, err := h.tm.Get("ssg", "new-content")
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

	w.WriteHeader(statusCode)
	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}
