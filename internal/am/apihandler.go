package am

import (
	"net/http"

	"github.com/google/uuid"
)

type APIHandler struct {
	*Handler
}

func NewAPIHandler(name string, opts ...Option) *APIHandler {
	handler := NewHandler(name, opts...)
	return &APIHandler{
		Handler: handler,
	}
}

func (h *APIHandler) Err(w http.ResponseWriter, code int, message string, err error) {
	var details string
	if err != nil {
		details = err.Error()
	}
	h.Log().Errorf("%s: %s (details: %s)", h.Name(), message, details)
	Respond(w, code, NewErrorResponse(message, ErrorCodeInternalError, details))
}

func (h *APIHandler) OK(w http.ResponseWriter, message string, data interface{}) {
	Respond(w, http.StatusOK, NewSuccessResponse(message, data))
}

func (h *APIHandler) Created(w http.ResponseWriter, message string, data interface{}) {
	Respond(w, http.StatusCreated, NewSuccessResponse(message, data))
}

func (h *APIHandler) ID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	id, err := PathID(r, "id")
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid ID in URL", err)
		return uuid.Nil, err
	}
	return id, nil
}
