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
	user := auth.NewUser("fakeuser", "Fake User")
	user.BaseModel.SetID(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	user.IsActive = true
	return user
}
