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
	h.renderContentForm(w, r, form, NewContent("", ""), "", http.StatusOK)
}

func (h *WebHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create content")
	h.Log().Infof("HX-Request header: %s", r.Header.Get("HX-Request"))
	ctx := r.Context()

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.renderContentForm(w, r, form, NewContent("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	err = form.Validate()
	if err != nil || form.HasErrors() {
		h.renderContentForm(w, r, form, NewContent("", ""), "Validation failed", http.StatusBadRequest)
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

	// If HTMX request, set HX-Redirect header and return
	if am.IsHTMXRequest(r) {
		redirectURL := am.EditPath(ssgPath, contentPath, content.ID())
		h.Log().Infof("HTMX request: Setting HX-Redirect to %s", redirectURL)
		w.Header().Set("HX-Redirect", redirectURL)
		h.renderContentForm(w, r, form, content, "", http.StatusOK)
		return
	}

	h.FlashInfo(w, r, "Content created")
	h.Redir(w, r, am.EditPath(ssgPath, contentPath, content.ID()), http.StatusSeeOther)
}

func (h *WebHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update content")
	ctx := r.Context()

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.renderContentForm(w, r, form, NewContent("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	err = form.Validate()
	if err != nil || form.HasErrors() {
		h.renderContentForm(w, r, form, NewContent("", ""), "Validation failed", http.StatusBadRequest)
		return
	}

	content := ToContentFromForm(form)
	content.GenUpdateValues()

	user := h.sampleUserInSession(r)
	content.UserID = user.ID()

	err = h.service.UpdateContent(ctx, content)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	if am.IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		// Return a small HTML fragment with the current timestamp for auto-save feedback
		_, _ = w.Write([]byte("<div id=\"save-status\" data-timestamp=\"" + am.Now().Format(am.TimeFormat) + "\"></div>"))
		return
	}

	h.FlashInfo(w, r, "Content updated")
	h.Redir(w, r, am.EditPath(ssgPath, contentPath, content.ID()), http.StatusSeeOther)
}

func (h *WebHandler) renderContentForm(w http.ResponseWriter, r *http.Request, form ContentForm, content Content, errorMessage string, statusCode int) {
	h.Log().Info("Render content form")
	h.Log().Infof("renderContentForm - form: %+v", form)
	h.Log().Infof("renderContentForm - form.ID: %s", form.ID)
	h.Log().Infof("renderContentForm - content: %+v", content)
	if content.BaseModel != nil {
		h.Log().Infof("renderContentForm - content.BaseModel ID: %s", content.BaseModel.ID())
	} else {
		h.Log().Info("renderContentForm - content.BaseModel is nil")
	}
	ctx := r.Context()

	sections, err := h.service.GetSections(ctx) // NOTE: This value should be cached
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, content)
	page.SetForm(form)

	if content.IsZero() {
		page.Form.SetAction(am.CreatePath(ssgPath, contentPath))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Form.SetAction(am.UpdatePath(ssgPath, contentPath))
		page.Form.SetSubmitButtonText("Update")
	}

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

func (h *WebHandler) EditContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit content")
	ctx := r.Context()

	id := r.URL.Query().Get("id")
	h.Log().Infof("EditContent: Received ID: %s", id)
	if id == "" {
		h.Err(w, nil, am.ErrBadRequest, http.StatusBadRequest)
		return
	}

	content, err := h.service.GetContent(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResource, http.StatusInternalServerError)
		return
	}

	form := ToContentForm(r, content)
	h.renderContentForm(w, r, form, content, "", http.StatusOK)
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
