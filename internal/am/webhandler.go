package am

import (
	"bytes"
	"net/http"
)

type WebHandler struct {
	*Handler
	tm    *TemplateManager
	flash *FlashManager
}

func NewWebHandler(tm *TemplateManager, flash *FlashManager, options ...Option) *WebHandler {
	handler := NewHandler("web-handler", options...)
	return &WebHandler{
		Handler: handler,
		tm:      tm,
		flash:   flash,
	}
}

func (h *WebHandler) Tmpl() *TemplateManager {
	return h.tm
}

func (h *WebHandler) Flash() *FlashManager {
	return h.flash
}

func (h *WebHandler) Success(w http.ResponseWriter, r *http.Request, msg string) {
	h.Flash().AddFlash(r, NotificationType.Success, msg)
}

func (h *WebHandler) Info(w http.ResponseWriter, r *http.Request, msg string) {
	h.Flash().AddFlash(r, NotificationType.Info, msg)
}

func (h *WebHandler) Warn(w http.ResponseWriter, r *http.Request, msg string) {
	h.Flash().AddFlash(r, NotificationType.Warn, msg)
}

func (h *WebHandler) Error(w http.ResponseWriter, r *http.Request, msg string) {
	h.Flash().AddFlash(r, NotificationType.Error, msg)
}

func (h *WebHandler) Debug(w http.ResponseWriter, r *http.Request, msg string) {
	h.Flash().AddFlash(r, NotificationType.Debug, msg)
}

func (h *WebHandler) OK(w http.ResponseWriter, r *http.Request, buf *bytes.Buffer, statusCode int) {
	h.Flash().ClearFlashCookie(w)

	w.WriteHeader(statusCode)
	_, err := buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) Redir(w http.ResponseWriter, r *http.Request, path string, status int) {
	flash := h.Flash().GetFlash(r)

	if len(flash.Notifications) > 0 {
		h.Flash().SetFlashInCookie(w, flash)
	} else {
		h.Flash().ClearFlashCookie(w)
	}

	http.Redirect(w, r, path, status)
}
