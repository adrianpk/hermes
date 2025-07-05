package ssg

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/adrianpk/hermes/internal/feat/auth"
	"github.com/google/uuid"
)

const (
	ssgPath = "/ssg"
)

type WebHandler struct {
	*am.WebHandler
	service Service
}

func NewWebHandler(tm *am.TemplateManager, flash *am.FlashManager, service Service, options ...am.Option) *WebHandler {
	handler := am.NewWebHandler(tm, flash, options...)
	return &WebHandler{
		WebHandler: handler,
		service:    service,
	}
}

// sampleUserInSession returns a fake user for now.
func (h *WebHandler) sampleUserInSession(r *http.Request) auth.User {
	return auth.User{
		BaseModel: am.NewModel(
			am.WithID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			am.WithType("user"),
		),
		Username: "fakeuser",
		Name:     "Fake User",
		IsActive: true,
	}
}
