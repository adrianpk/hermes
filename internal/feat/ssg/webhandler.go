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

	content := NewContent("", "")

	page := am.NewPage(r, content)
	page.SetFormAction(am.CreatePath(ssgPath, contentPath))
	page.SetFormButtonText("Create")

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

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}
