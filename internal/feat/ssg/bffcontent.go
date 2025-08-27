package ssg

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

func (h *BFF) NewContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New content form")
	form := NewContentForm(r)
	h.renderContentForm(w, r, form, NewContent("", ""), "", http.StatusOK)
}

func (h *BFF) CreateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create content")

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.renderContentForm(w, r, form, NewContent("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil || form.HasErrors() {
		h.renderContentForm(w, r, form, NewContent("", ""), "Validation failed", http.StatusBadRequest)
		return
	}

	content := ToContent(form)
	user := h.sampleUserInSession(r)
	content.UserID = user.GetID()

	var createdContent Content
	err = h.apiClient.Post("/contents", content, &createdContent)
	if err != nil {
		h.Err(w, err, "Failed to create content via API", http.StatusInternalServerError)
		return
	}

	if am.IsHTMXRequest(r) {
		redirectURL := am.EditPath(ssgPath, contentPath, createdContent.GetID())
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.FlashInfo(w, r, "Content created")
	h.Redir(w, r, am.EditPath(ssgPath, contentPath, createdContent.GetID()), http.StatusSeeOther)
}

func (h *BFF) EditContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit content")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing content ID", http.StatusBadRequest)
		return
	}

	var content Content
	path := fmt.Sprintf("/contents/%s", idStr)
	err := h.apiClient.Get(path, &content)
	if err != nil {
		h.Err(w, err, "Failed to get content from API", http.StatusInternalServerError)
		return
	}

	form := ToContentForm(r, content)
	h.renderContentForm(w, r, form, content, "", http.StatusOK)
}

func (h *BFF) UpdateContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update content")

	form, err := ContentFormFromRequest(r)
	if err != nil {
		h.renderContentForm(w, r, form, NewContent("", ""), "Invalid form data", http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil || form.HasErrors() {
		h.renderContentForm(w, r, form, NewContent("", ""), "Validation failed", http.StatusBadRequest)
		return
	}

	content := ToContent(form)
	user := h.sampleUserInSession(r)
	content.SetUpdatedBy(user.GetID())

	path := fmt.Sprintf("/contents/%s", content.GetID())
	err = h.apiClient.Put(path, content, nil)
	if err != nil {
		h.Err(w, err, "Failed to update content via API", http.StatusInternalServerError)
		return
	}

	if am.IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<div id=\"save-status\" data-timestamp=\"" + am.Now().Format(am.TimeFormat) + "\"></div>"))
		return
	}

	h.FlashInfo(w, r, "Content updated successfully")
	h.Redir(w, r, am.EditPath(ssgPath, contentPath, content.GetID()), http.StatusSeeOther)
}

func (h *BFF) ListContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List content")

	var contents []Content
	err := h.apiClient.Get("/contents", &contents)
	if err != nil {
		h.Err(w, err, "Failed to get contents from API", http.StatusInternalServerError)
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

func (h *BFF) ShowContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show content")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing content ID", http.StatusBadRequest)
		return
	}

	var content Content
	path := fmt.Sprintf("/contents/%s", idStr)
	err := h.apiClient.Get(path, &content)
	if err != nil {
		h.Err(w, err, "Failed to get content from API", http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, content)
	page.Name = "Show Content"

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(content, "Back")

	tmpl, err := h.Tmpl().Get(ssgFeat, "show-content")
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

func (h *BFF) DeleteContent(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete content")

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idStr := r.Form.Get("id")
	if idStr == "" {
		h.Err(w, nil, "Missing content ID", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("/contents/%s", idStr)
	err := h.apiClient.Delete(path)
	if err != nil {
		h.Err(w, err, "Failed to delete content via API", http.StatusInternalServerError)
		return
	}

	h.FlashInfo(w, r, "Content deleted successfully")
	h.Redir(w, r, am.ListPath(ssgPath, contentPath), http.StatusSeeOther)
}

func (h *BFF) renderContentForm(w http.ResponseWriter, r *http.Request, form ContentForm, content Content, errorMessage string, statusCode int) {
	var sections []Section
	err := h.apiClient.Get("/sections", &sections)
	if err != nil {
		h.Err(w, err, "Failed to get sections from API", http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, content)
	page.SetForm(&form)
	page.AddSelect("sections", am.ToSelectOpt(sections))

	if content.IsZero() {
		page.Name = "New Content"
		page.IsNew = true
		page.Form.SetAction(am.CreatePath(ssgPath, contentPath))
		page.Form.SetSubmitButtonText("Create")
	} else {
		page.Name = "Edit Content"
		page.IsNew = false
		page.Form.SetAction(am.UpdatePath(ssgPath, contentPath))
		page.Form.SetSubmitButtonText("Update")
	}

	menu := page.NewMenu(ssgPath)
	menu.AddListItem(content, "Back")

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
