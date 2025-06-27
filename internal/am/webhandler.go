package am

import (
	"bytes"
	"net/http"
)

type WebHandler struct {
	*Handler
	tm *TemplateManager
	fm *FlashManager
}

func NewWebHandler(tm *TemplateManager, fm *FlashManager, options ...Option) *WebHandler {
	handler := NewHandler("web-handler", options...)
	return &WebHandler{
		Handler: handler,
		tm:      tm,
		fm:      fm,
	}
}

func (h *WebHandler) Tmpl() *TemplateManager {
	return h.tm
}

func (h *WebHandler) FlashManager() *FlashManager {
	return h.fm
}

func (h *WebHandler) GetFlash(r *http.Request) Flash {
	return h.FlashManager().GetFlash(r)
}

func (h *WebHandler) FlashSuccess(w http.ResponseWriter, r *http.Request, msg string) {
	h.FlashManager().AddFlash(r, NotificationType.Success, msg)
}

func (h *WebHandler) FlashInfo(w http.ResponseWriter, r *http.Request, msg string) {
	h.FlashManager().AddFlash(r, NotificationType.Info, msg)
}

func (h *WebHandler) FlashWarn(w http.ResponseWriter, r *http.Request, msg string) {
	h.FlashManager().AddFlash(r, NotificationType.Warn, msg)
}

func (h *WebHandler) FlashError(w http.ResponseWriter, r *http.Request, msg string) {
	h.FlashManager().AddFlash(r, NotificationType.Error, msg)
}

func (h *WebHandler) FlashDebug(w http.ResponseWriter, r *http.Request, msg string) {
	h.FlashManager().AddFlash(r, NotificationType.Debug, msg)
}

func (h *WebHandler) OK(w http.ResponseWriter, r *http.Request, buf *bytes.Buffer, statusCode int) {
	h.FlashManager().ClearFlashCookie(w)

	w.WriteHeader(statusCode)
	_, err := buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) Redir(w http.ResponseWriter, r *http.Request, path string, status int) {
	flash := h.FlashManager().GetFlash(r)

	if len(flash.Notifications) > 0 {
		h.FlashManager().SetFlashInCookie(w, flash)
	} else {
		h.FlashManager().ClearFlashCookie(w)
	}

	http.Redirect(w, r, path, status)
}
