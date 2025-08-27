package am

import (
	"fmt"
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
	idStr := r.PathValue("id")
	if idStr == "" {
		h.Err(w, http.StatusBadRequest, "ID is missing in the URL", fmt.Errorf("ID is missing in the URL"))
		return uuid.Nil, fmt.Errorf("ID is missing in the URL")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, http.StatusBadRequest, "Invalid ID format", fmt.Errorf("invalid ID format: %w", err))
		return uuid.Nil, fmt.Errorf("invalid ID format: %w", err)
	}

	return id, nil
}
